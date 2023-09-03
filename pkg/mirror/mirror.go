package mirror

import (
	"errors"
	"net/url"

	"github.com/coffeebeats/gdenv/pkg/godot"
)

const filenameChecksums = "SHA512-SUMS.txt"

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
	name string
	url  *url.URL
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
