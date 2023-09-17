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
//
// NOTE: A dummy method is defined on this interface in order to (i) restrict
// outside implementers and (ii) ensure the correct 'Archive' and 'Artifact'
// types are used during checksum extraction.
type Checksums[T artifact.Artifact, U archive.Archive[T]] interface {
	artifact.Artifact
	artifact.Versioned

	supports(U) // A private dummy method to restrict implementers.
}

/* -------------------------------------------------------------------------- */
/*                              Struct: checksums                             */
/* -------------------------------------------------------------------------- */

// A shared implementation of a checksums file 'Artifact'; this should be
// wrapped by user-facing types.
type checksums struct {
	version version.Version
}
