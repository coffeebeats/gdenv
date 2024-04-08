package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"

	"github.com/coffeebeats/gdenv/pkg/progress"
)

const (
	// Configure common retry policies for clients.
	retryCount   = 3
	retryWait    = time.Second
	retryWaitMax = time.Second
)

var (
	ErrClientConfiguration    = errors.New("client misconfigured")
	ErrHTTPResponseStatusCode = errors.New("received error status code")
	ErrInvalidURL             = errors.New("invalid URL")
	ErrMissingSize            = errors.New("missing progress size")
	ErrMissingURL             = errors.New("missing URL")
	ErrRequestFailed          = errors.New("request failed")
	ErrUnexpectedRedirect     = errors.New("unexpected redirect")
)

type progressKey struct{}

/* -------------------------------------------------------------------------- */
/*                             Function: ParseURL                             */
/* -------------------------------------------------------------------------- */

// Returns a parsed URL or fails if it's invalid.
func ParseURL(urlBaseRaw string, urlPartsRaw ...string) (*url.URL, error) {
	if urlBaseRaw == "" {
		return nil, ErrMissingURL
	}

	urlRaw, err := url.JoinPath(urlBaseRaw, urlPartsRaw...)
	if err != nil {
		return nil, errors.Join(ErrInvalidURL, err)
	}

	// NOTE: Use the stricter 'ParseRequestURI' function instead of 'Parse'.
	urlParsed, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		return nil, errors.Join(ErrInvalidURL, err)
	}

	if urlParsed.Host == "" || urlParsed.Scheme == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, urlRaw)
	}

	return urlParsed, nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithProgress creates a sub-context with an associated progress reporter. The
// result can be passed to file download functions in this package to get
// updates on download progress.
func WithProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKey{}, p)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Client                               */
/* -------------------------------------------------------------------------- */

// A struct implementing an HTTP client with simple methods for file downloads.
type Client struct {
	restyClient *resty.Client
}

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Client' with default settings for mirrors.
func New() *Client {
	restyClient := resty.New()

	restyClient.SetRetryCount(retryCount)
	restyClient.SetRetryWaitTime(retryWait)
	restyClient.SetRetryMaxWaitTime(retryWaitMax)

	// Disable redirects by default.
	restyClient.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Retry on any request execution error or retryable HTTP status code.
	restyClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			s := r.StatusCode()

			return err != nil ||
				s == http.StatusRequestTimeout || // 408
				s == http.StatusTooManyRequests || // 429
				s == http.StatusInternalServerError || // 500
				s == http.StatusBadGateway || // 502
				s == http.StatusServiceUnavailable || // 503
				s == http.StatusGatewayTimeout // 504
		},
	)

	// Add logging when a retry occurs.
	restyClient.AddRetryHook(func(r *resty.Response, err error) {
		if r.Request != nil {
			select {
			case <-r.Request.Context().Done():
				return
			default:
			}
		}

		log.Warn("Retrying request due to error:", err, fmt.Sprintf("(%s)", r.Status()))
	})

	// Disable internal resty client logging.
	restyClient.SetLogger(silentLogger{})

	return &Client{restyClient}
}

/* -------------------- Function: NewWithRedirectDomains -------------------- */

// Creates a new 'Client' with the provided domains allowed for redirects. If
// none are provided, then no redirects are permitted.
func NewWithRedirectDomains(domains ...string) *Client {
	var p resty.RedirectPolicy

	switch len(domains) {
	case 0:
		p = resty.NoRedirectPolicy()
	default:
		p = resty.DomainCheckRedirectPolicy(domains...)
	}

	client := New()

	client.restyClient.SetRedirectPolicy(p)

	return client
}

/* ----------------------------- Method: Exists ----------------------------- */

// Issues a 'HEAD' request to test whether or not the URL is reachable.
func (c *Client) Exists(ctx context.Context, urlBaseRaw string, urlPartsRaw ...string) (bool, error) {
	if urlBaseRaw == "" {
		return false, ErrMissingURL
	}

	urlRaw, err := url.JoinPath(urlBaseRaw, urlPartsRaw...)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	// NOTE: Use the stricter 'ParseRequestURI' function instead of 'Parse'.
	urlParsed, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		return false, errors.Join(ErrInvalidURL, err)
	}

	if urlParsed.Host == "" || urlParsed.Scheme == "" {
		return false, fmt.Errorf("%w: %s", ErrInvalidURL, urlParsed.String())
	}

	// Use a no-op response handler, as just the response code is used.
	err = c.head(ctx, urlParsed, func(r *resty.Response) error {
		// Redirects should be followed by the client, not accepted as a valid
		// result for 'Exists'. Return an error so the caller knows the client
		// is incorrectly configured.
		if r.StatusCode() >= http.StatusMultipleChoices && r.StatusCode() < http.StatusBadRequest {
			return errors.Join(ErrClientConfiguration, ErrUnexpectedRedirect)
		}

		return nil
	})

	switch {
	// A response error occurred, indicating there's a problem reaching the URL.
	case errors.Is(err, ErrHTTPResponseStatusCode) || errors.Is(err, ErrRequestFailed):
		return false, nil
	// A request execution error occurred.
	case err != nil:
		return false, err

	default:
		return true, nil
	}
}

/* ---------------------------- Method: Download ---------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers. Reports progress to a 'progress.Progress' set on the
// provided context.
func (c *Client) Download(ctx context.Context, u *url.URL, w ...io.Writer) error {
	return c.get(ctx, u, func(r *resty.Response) error {
		if r.RawResponse.ContentLength > 0 { // No progress to report if '0'.
			// Report progress if set on the context.
			if p, ok := ctx.Value(progressKey{}).(*progress.Progress); ok && p != nil {
				if err := p.SetTotal(uint64(r.RawResponse.ContentLength)); err != nil {
					return err
				}

				w = append(w, progress.NewWriter(p))
			}
		}

		// Copy the asset contents into provided writers.
		_, err := io.Copy(io.MultiWriter(w...), r.RawBody())

		return err
	})
}

/* --------------------------- Method: DownloadTo --------------------------- */

// Downloads the provided asset to a specified file 'out'. Reports progress to
// a 'progress.Progress' set on the provided context.
func (c *Client) DownloadTo(ctx context.Context, u *url.URL, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	return c.get(ctx, u, func(r *resty.Response) error {
		var w io.Writer = f

		if r.RawResponse.ContentLength > 0 { // No progress to report if '0'.
			// Report progress if set on the context.
			if p, ok := ctx.Value(progressKey{}).(*progress.Progress); ok && p != nil {
				if err := p.SetTotal(uint64(r.RawResponse.ContentLength)); err != nil {
					return err
				}

				w = io.MultiWriter(f, progress.NewWriter(p))
			}
		}

		// Copy the response contents into the writer.
		_, err := io.Copy(w, r.RawBody())

		return err
	})
}

/* --------------------------- Method: RestyClient -------------------------- */

// RestyClient returns the underlying 'resty.Client'. Exposing this allows for
// tests in other packages to stub responses.
//
// TODO: Find an alternative way of allowing other packages to stub responses
// without exposing the use of 'resty'.
func (c *Client) RestyClient() *resty.Client {
	return c.restyClient
}

/* ------------------------------- Method: get ------------------------------ */

// A convenience wrapper around execute which issues a 'GET' request.
func (c *Client) get(ctx context.Context, u *url.URL, h func(*resty.Response) error) error {
	return execute(ctx, c.restyClient.R(), resty.MethodGet, u.String(), h)
}

/* ------------------------------ Method: head ------------------------------ */

// A convenience wrapper around execute which issues a 'HEAD' request.
func (c *Client) head(ctx context.Context, u *url.URL, h func(*resty.Response) error) error {
	return execute(ctx, c.restyClient.R(), resty.MethodHead, u.String(), h)
}

/* ---------------------------- Function: execute --------------------------- */

// Executes the provided request, but delegates the response handling to the
// provided function. The handler should *not* close the response, as that's
// handled by this function.
func execute(ctx context.Context, req *resty.Request, m, u string, h func(*resty.Response) error) error {
	res, err := req.
		SetContext(ctx).             // Allow canceling the request.
		SetDoNotParseResponse(true). // Take over response parsing (requires manually closing response body).
		Execute(m, u)
	if err != nil {
		return errors.Join(ErrRequestFailed, err)
	}

	defer res.RawBody().Close()

	if res.IsError() {
		return fmt.Errorf("%w: %w: %d", ErrRequestFailed, ErrHTTPResponseStatusCode, res.StatusCode())
	}

	return h(res)
}

/* -------------------------------------------------------------------------- */
/*                             Type: silentLogger                             */
/* -------------------------------------------------------------------------- */

// silentLogger is a 'resty.Logger' implementation that emits no logs.
type silentLogger struct{}

// Compile-time verification that 'silentLogger' implements 'resty.Logger'.
var _ resty.Logger = (*silentLogger)(nil)

func (l silentLogger) Debugf(string, ...interface{}) {}
func (l silentLogger) Errorf(string, ...interface{}) {}
func (l silentLogger) Warnf(string, ...interface{})  {}
