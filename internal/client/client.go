package client

import (
	"errors"
	"io"
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
	ErrMissingSize = errors.New("missing progress size")
	ErrNetwork     = errors.New("network request failure")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Client                               */
/* -------------------------------------------------------------------------- */

// A struct implementing an HTTP client with simple methods for file downloads.
type Client struct {
	restyClient *resty.Client
}

// Validate at compile-time that 'Client' implements 'FileDownloader'.
var _ FileDownloader[*url.URL] = &Client{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Client' with default settings for mirrors.
func New() Client {
	restyClient := resty.New()

	restyClient.SetRetryCount(retryCount)
	restyClient.SetRetryWaitTime(retryWait)
	restyClient.SetRetryMaxWaitTime(retryWaitMax)

	// Disable redirects by default.
	restyClient.SetRedirectPolicy(resty.NoRedirectPolicy())

	// Retry on any error response.
	restyClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return err != nil || r.IsError()
		},
	)

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

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (c Client) Download(u *url.URL, w ...io.Writer) error {
	return get(c, u, func(r *resty.Response) error {
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

	return get(c, u, func(r *resty.Response) error {
		// Copy the response contents into the writer.
		_, err := io.Copy(f, r.RawBody())

		return err
	})
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (c Client) DownloadToWithProgress(u *url.URL, out string, p *progress.Progress) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	return get(c, u, func(r *resty.Response) error {
		w := progress.NewWriter(p)

		// Copy the asset contents into the writer.
		_, err := io.Copy(io.MultiWriter(f, &w), r.RawBody())

		return err
	})
}

/* ------------------------------ Function: get ----------------------------- */

func get(c Client, u *url.URL, h func(*resty.Response) error) error {
	req := c.restyClient.R()

	// Assume control of response parsing.
	req.SetDoNotParseResponse(true)

	// Issue the HTTP request.
	res, err := req.Get(u.String())
	if err != nil {
		return errors.Join(ErrNetwork, err)
	}

	defer res.RawBody().Close()

	return h(res)
}
