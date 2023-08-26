package mirror

import (
	"io"
	"net/url"
	"time"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

// Configure common retry policies for mirrors.
const (
	retryCount   = 3
	retryWait    = time.Second
	retryWaitMax = 10 * time.Second
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Asset                              */
/* -------------------------------------------------------------------------- */

// An interface representing a file/directory that can be downloaded from a
// Godot project mirror.
type Asset interface {
	Name() string
	URL() url.URL

	Download(w io.Writer) error
}

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// An interface representing a host of Godot project artifacts.
type Mirror interface {
	Checksum(v godot.Version) (Asset, error)
	Executable(ex godot.Executable) (Asset, error)

	Has(ex godot.Executable) bool
}

/* -------------------------------------------------------------------------- */
/*                                Struct: asset                               */
/* -------------------------------------------------------------------------- */

// The GitHub-specific implementation of a release asset.
type asset struct {
	client *resty.Client
	name   string
	url    url.URL
}

/* ------------------------------ Method: Name ------------------------------ */

func (a *asset) Name() string {
	return a.name
}

/* ------------------------------- Method: URL ------------------------------ */

func (a *asset) URL() url.URL {
	return a.url
}

/* ---------------------------- Method: Download ---------------------------- */

func (a *asset) Download(w io.Writer) error {
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
