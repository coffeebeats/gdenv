package mirror

import (
	"errors"

	"github.com/coffeebeats/gdenv/internal/version"
	"github.com/coffeebeats/gdenv/pkg/godot"
)

var (
	ErrInvalidSpecification = errors.New("invalid specification")
	ErrInvalidURL           = errors.New("invalid URL")
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// An interface specifying methods for retrieving information about assets
// available for download via a mirror host.
type Mirror interface {
	Checksum(v version.Version) (Asset, error)
	Executable(ex godot.Executable) (Asset, error)

	// Issues a request to see if the mirror host has the specific version.
	Has(v version.Version) bool

	// Checks whether the version is broadly supported by the mirror. No network
	// request is issued, but this does not guarantee the host has the version.
	// To check whether the host has the version definitively via the network,
	// use the 'Has' method.
	Supports(v version.Version) bool
}
