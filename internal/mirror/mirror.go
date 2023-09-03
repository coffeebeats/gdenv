package mirror

import (
	"errors"

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
	Checksum(v godot.Version) (Asset, error)
	Executable(ex godot.Executable) (Asset, error)

	Supports(v godot.Version) bool
}
