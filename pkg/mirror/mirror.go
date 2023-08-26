package mirror

import (
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

// Configure common retry policies for mirrors.
const (
	filenameChecksums = "SHA512-SUMS.txt"

	retryCount   = 3
	retryWait    = time.Second
	retryWaitMax = 10 * time.Second
)

var (
	ErrInvalidSpecification = errors.New("invalid specification")
	ErrInvalidURL           = errors.New("invalid URL")
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

func (a Asset) Name() string {
	return a.name
}

/* ------------------------------- Method: URL ------------------------------ */

func (a Asset) URL() *url.URL {
	return a.url
}

/* ---------------------------- Method: Download ---------------------------- */

func (a Asset) Download(w io.Writer) error {
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
