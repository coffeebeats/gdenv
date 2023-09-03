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
	ErrFileSystem      = errors.New("file system operation failed")
	ErrIO              = errors.New("I/O operation failed")
	ErrMissingResponse = errors.New("missing response")
	ErrMissingClient   = errors.New("missing client")
	ErrMissingSize     = errors.New("missing progress size")
	ErrNetwork         = errors.New("network request failure")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Client                               */
/* -------------------------------------------------------------------------- */

// A struct implementing an HTTP client with simple methods for file downloads.
type Client struct {
	client *resty.Client
}

// Validate at compile-time that 'Client' implements 'FileDownloader'.
var _ FileDownloader = &Client{} //nolint:exhaustruct

/* ---------------------------- Function: Default --------------------------- */

// Configures a default 'Client' for mirrors.
func Default() Client {
	client := resty.New()

	client.SetRetryCount(retryCount)
	client.SetRetryWaitTime(retryWait)
	client.SetRetryMaxWaitTime(retryWaitMax)

	// Retry on any error response.
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return err != nil || r.IsError()
		},
	)

	return Client{client}
}

/* ------------------------ Method: AllowRedirectsTo ------------------------ */

// Modifies the client's redirect policies to only allow redirects to the
// provided domains. If none are provided, then no redirects are permitted.
func (c *Client) AllowRedirectsTo(d ...string) {
	var p resty.RedirectPolicy

	switch len(d) {
	case 0:
		p = resty.NoRedirectPolicy()
	default:
		p = resty.DomainCheckRedirectPolicy(d...)
	}

	c.client.SetRedirectPolicy(p)
}

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (c *Client) Download(u *url.URL, w ...io.Writer) error {
	return get(c, u, func(r *resty.Response) error {
		if r == nil {
			return ErrMissingResponse
		}

		// Copy the asset contents into provided writers.
		if _, err := io.Copy(io.MultiWriter(w...), r.RawBody()); err != nil {
			return errors.Join(ErrIO, err)
		}

		return nil
	})
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (c *Client) DownloadTo(u *url.URL, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return errors.Join(ErrFileSystem, err)
	}

	defer f.Close()

	return get(c, u, func(r *resty.Response) error {
		if r == nil {
			return ErrMissingResponse
		}

		// Copy the response contents into the writer.
		if _, err := io.Copy(f, r.RawBody()); err != nil {
			return errors.Join(ErrIO, err)
		}

		return nil
	})
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (c *Client) DownloadToWithProgress(u *url.URL, out string, p *progress.Progress) error {
	f, err := os.Create(out)
	if err != nil {
		return errors.Join(ErrFileSystem, err)
	}

	defer f.Close()

	return get(c, u, func(r *resty.Response) error {
		if r == nil {
			return ErrMissingResponse
		}

		w := progress.NewWriter(p)

		// Copy the asset contents into the writer.
		if _, err := io.Copy(io.MultiWriter(f, w), r.RawBody()); err != nil {
			return errors.Join(ErrIO, err)
		}

		return nil
	})
}

/* ------------------------------ Function: get ----------------------------- */

func get(c *Client, u *url.URL, h func(*resty.Response) error) error {
	if c.client == nil {
		return ErrMissingClient
	}

	req := c.client.R()

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
