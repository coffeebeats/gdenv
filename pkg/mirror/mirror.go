package mirror

import (
	"io"
	"net/url"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
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
