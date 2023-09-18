package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/coffeebeats/gdenv/internal/progress"
	"github.com/go-resty/resty/v2"
)

const (
	// Configure common retry policies for clients.
	retryCount   = 3
	retryWait    = time.Second
	retryWaitMax = 10 * time.Second
)

var (
	ErrClientConfiguration    = errors.New("client misconfigured")
	ErrHTTPResponseStatusCode = errors.New("received error status code")
	ErrInvalidURL             = errors.New("invalid URL")
	ErrMissingSize            = errors.New("missing progress size")
	ErrRequestFailed          = errors.New("request failed")
	ErrUnexpectedRedirect     = errors.New("unexpected redirect")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Client                               */
/* -------------------------------------------------------------------------- */

// A struct implementing an HTTP client with simple methods for file downloads.
type Client struct {
	restyClient *resty.Client
}

// Validate at compile-time that 'Client' implements 'Downloader'.
var _ Downloader[*url.URL] = &Client{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Client' with default settings for mirrors.
func New() Client {
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
		log.Println("Retrying due to:", err, fmt.Sprintf("(%s)", r.Status()))
	})

	return Client{restyClient}
}

/* -------------------- Function: NewWithRedirectDomains -------------------- */

// Creates a new 'Client' with the provided domains allowed for redirects. If
// none are provided, then no redirects are permitted.
func NewWithRedirectDomains(domains ...string) Client {
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
func (c Client) Exists(urlRaw string) (bool, error) {
	// NOTE: Use the stricter 'ParseRequestURI' function instead of 'Parse'.
	urlParsed, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		return false, errors.Join(ErrInvalidURL, err)
	}

	if urlParsed.Host == "" || urlParsed.Scheme == "" {
		return false, fmt.Errorf("%w: %s", ErrInvalidURL, urlRaw)
	}

	// Use a no-op response handler, as just the response code is used.
	err = c.head(urlParsed, func(r *resty.Response) error {
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
	case errors.Is(err, ErrHTTPResponseStatusCode):
		return false, nil
	// A request execution error occurred.
	case err != nil:
		return false, err

	default:
		return true, nil
	}
}

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (c Client) Download(u *url.URL, w ...io.Writer) error {
	return c.get(u, func(r *resty.Response) error {
		// Copy the asset contents into provided writers.
		_, err := io.Copy(io.MultiWriter(w...), r.RawBody())

		return err
	})
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (c Client) DownloadTo(u *url.URL, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	return c.get(u, func(r *resty.Response) error {
		// Copy the response contents into the writer.
		_, err := io.Copy(f, r.RawBody())

		return err
	})
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
//
// NOTE: The provided 'Progress' struct will be reconfigured as needed.
func (c Client) DownloadToWithProgress(u *url.URL, out string, p *progress.Progress) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	return c.get(u, func(r *resty.Response) error {
		// Reset any pre-existing progress in the 'Progress' reporter.
		p.Reset()

		// Set the 'Progress' total based on the response header.
		if err := p.Total(uint64(r.RawResponse.ContentLength)); err != nil {
			return err
		}

		// Copy the asset contents into the file and progress writer.
		_, err := io.Copy(io.MultiWriter(f, progress.NewWriter(p)), r.RawBody())

		return err
	})
}

/* ------------------------------- Method: get ------------------------------ */

// A convenience wrapper around execute which issues a 'GET' request.
func (c Client) get(u *url.URL, h func(*resty.Response) error) error {
	return execute(c.restyClient.R(), resty.MethodGet, u.String(), h)
}

/* ------------------------------ Method: head ------------------------------ */

// A convenience wrapper around execute which issues a 'HEAD' request.
func (c Client) head(u *url.URL, h func(*resty.Response) error) error {
	return execute(c.restyClient.R(), resty.MethodHead, u.String(), h)
}

/* ---------------------------- Function: execute --------------------------- */

// Executes the provided request, but delegates the response handling to the
// provided function. The handler should *not* close the response, as that's
// handled by this function.
func execute(req *resty.Request, m, u string, h func(*resty.Response) error) error {
	// Take over response parsing (requires response body to be manually closed).
	req.SetDoNotParseResponse(true)

	res, err := req.Execute(m, u)
	if err != nil {
		return errors.Join(ErrRequestFailed, err)
	}

	defer res.RawBody().Close()

	if res.IsError() {
		return fmt.Errorf("%w: %w: %d", ErrRequestFailed, ErrHTTPResponseStatusCode, res.StatusCode())
	}

	return h(res)
}
