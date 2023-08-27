package mirror

import (
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

const (
	filenameChecksums = "SHA512-SUMS.txt"

	// Configure common retry policies for mirrors.
	retryCount   = 3
	retryWait    = time.Second
	retryWaitMax = 10 * time.Second
)

var (
	ErrFileSystem           = errors.New("file system operation failed")
	ErrInvalidSpecification = errors.New("invalid specification")
	ErrInvalidURL           = errors.New("invalid URL")
	ErrMissingClient        = errors.New("missing client")
	ErrNetwork              = errors.New("network request failure")
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// An interface representing a host of Godot project artifacts.
type Mirror interface {
	Checksum(v godot.Version) (Asset, error)
	Executable(ex godot.Executable) (Asset, error)

	Supports(v godot.Version) bool
}

/* -------------------------------------------------------------------------- */
/*                                Struct: Asset                               */
/* -------------------------------------------------------------------------- */

// A struct representing a file/directory that can be downloaded from a Godot
// project mirror.
type Asset struct {
	client *resty.Client
	name   string
	url    *url.URL
}

/* ------------------------------ Method: Name ------------------------------ */

// Returns the filename of the asset to download.
func (a Asset) Name() string {
	return a.name
}

/* ------------------------------- Method: URL ------------------------------ */

// Returns the URL of the asset to download.
func (a Asset) URL() *url.URL {
	return a.url
}

/* ------------------------------ Method: Size ------------------------------ */

// Checks the size of the asset via the "Content-Length" response header.
func (a Asset) Size() (int64, error) {
	if a.client == nil {
		return 0, ErrMissingClient
	}

	// Issue the HTTP request.
	res, err := a.client.R().Head(a.url.String())
	if err != nil {
		return 0, errors.Join(ErrNetwork, err)
	}

	return res.RawResponse.ContentLength, nil
}

/* ---------------------------- Method: Download ---------------------------- */

// Downloads the asset, outputting the contents to the passed-in 'io.Writer'.
func (a Asset) Download(w io.Writer) error {
	if a.client == nil {
		return ErrMissingClient
	}

	// Assume control of response parsing.
	a.client.SetDoNotParseResponse(true)
	defer a.client.SetDoNotParseResponse(false)

	// Issue the HTTP request.
	res, err := a.client.R().Get(a.url.String())
	if err != nil {
		return errors.Join(ErrNetwork, err)
	}

	defer res.RawBody().Close()

	// Copy the asset contents into the writer.
	if _, err := io.Copy(w, res.RawBody()); err != nil {
		return errors.Join(ErrFileSystem, err)
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                             Function: newClient                            */
/* -------------------------------------------------------------------------- */

// Configures a default HTTP client for mirrors.
func newClient() *resty.Client {
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
