package checksum

import (
	"errors"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

var ErrChecksumsUnsupported = errors.New("version precedes checksums")

/* -------------------------------------------------------------------------- */
/*                            Interface: Checksums                            */
/* -------------------------------------------------------------------------- */

// An interface for an 'Artifact' representing a checksums file.
type Checksums[T artifact.Artifact, U archive.Archive[T]] interface {
	artifact.Artifact
	artifact.Versioned

	// NOTE: This dummy method is defined in order to (i) restrict outside
	// implementers and (ii) ensure the correct 'Archive' and 'Artifact' types
	// are used during checksum extraction.
	supports(U)
}

/* -------------------------------------------------------------------------- */
/*                              Struct: checksums                             */
/* -------------------------------------------------------------------------- */

// A shared implementation of a checksums file 'Artifact'; this should be
// wrapped by user-facing types.
type checksums struct {
	version version.Version
}
