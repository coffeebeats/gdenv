package mirror

import (
	"errors"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
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
	ExecutableArchive(ex executable.Executable) (artifact.Remote[executable.Archive], error)
	ExecutableArchiveChecksums(v version.Version) (artifact.Remote[checksum.Executable], error)

	SourceArchive(v version.Version) (artifact.Remote[source.Archive], error)
	SourceArchiveChecksums(v version.Version) (artifact.Remote[checksum.Source], error)

	// Issues a request to see if the mirror host has the specific version.
	CheckIfSupports(v version.Version) bool

	// Checks whether the version is broadly supported by the mirror. No network
	// request is issued, but this does not guarantee the host has the version.
	// To check whether the host has the version definitively via the network,
	// use the 'Has' method.
	Supports(v version.Version) bool
}
