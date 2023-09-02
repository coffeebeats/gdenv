package mirror

import (
	"errors"
	"io"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
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

type Client struct {
	client *resty.Client
}

/* ---------------------------- Method: Download ---------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (c *Client) Download(a Asset, w ...io.Writer) error {
	return get(c, a.URL(), func(r *resty.Response) error {
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

/* --------------------------- Method: DownloadTo --------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (c *Client) DownloadTo(a Asset, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return errors.Join(ErrFileSystem, err)
	}

	defer f.Close()

	return get(c, a.URL(), func(r *resty.Response) error {
		if r == nil {
			return ErrMissingResponse
		}

		// Copy the asset contents into the writer.
		if _, err := io.Copy(f, r.RawBody()); err != nil {
			return errors.Join(ErrIO, err)
		}

		return nil
	})
}

/* --------------------- Method: DownloadToWithProgress --------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (c *Client) DownloadToWithProgress(a Asset, out string, p *float64) error {
	f, err := os.Create(out)
	if err != nil {
		return errors.Join(ErrFileSystem, err)
	}

	defer f.Close()

	return get(c, a.URL(), func(r *resty.Response) error {
		if r == nil {
			return ErrMissingResponse
		}

		w, err := NewProgressWriter(r.RawResponse.ContentLength, p)
		if err != nil {
			return err
		}

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

/* -------------------------------------------------------------------------- */
/*                        Function: defaultRestyClient                        */
/* -------------------------------------------------------------------------- */

// Configures a default 'resty' (HTTP) client for mirrors.
func defaultRestyClient() *resty.Client {
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

	return client
}

/* -------------------------------------------------------------------------- */
/*                           Struct: ProgressWriter                           */
/* -------------------------------------------------------------------------- */

// An 'io.Writer' implementation that simply tracks the percentage of bytes
// written (i.e. progress).
type ProgressWriter struct {
	progress    *float64
	bytes, size int64
}

// Validate at compile-time that 'ProgressWriter' implements 'io.Writer'.
var _ io.Writer = &ProgressWriter{} //nolint:exhaustruct

/* ----------------------- Function: NewProgressWriter ---------------------- */

// Creates a new 'ProgressWriter' with the specified byte total and progress
// pointer. The progress pointer will be updated as bytes are written.
//
// NOTE: It's the caller's responsibility to ensure that the initial 'size'
// total is correct so that the computed progress value is accurate.
func NewProgressWriter(size int64, progress *float64) (*ProgressWriter, error) {
	if size == 0 {
		return nil, ErrMissingSize
	}

	return &ProgressWriter{progress: progress, size: size, bytes: 0}, nil
}

/* ----------------------------- Impl: io.Writer ---------------------------- */

func (pw *ProgressWriter) Write(data []byte) (int, error) {
	pw.bytes += int64(len(data))

	if pw.size > 0 && pw.progress != nil {
		*(pw.progress) = float64(pw.bytes) / float64(pw.size)
	}

	return len(data), nil
}
