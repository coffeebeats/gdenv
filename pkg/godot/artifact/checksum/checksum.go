package checksum

import (
	"context"
	"errors"
	"hash"

	"github.com/coffeebeats/gdenv/internal/checksumutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

// NOTE: 'ErrChecksumMismatch' is an alias of the error returned by the
// underlying checksum implementation so that 'errors.Is' checks continue to
// match. 'ErrChecksumsUnsupported' is specific to Godot's releases - it has no
// meaning at the format level - so it is defined here.
var (
	ErrChecksumMismatch     = checksumutil.ErrMismatch
	ErrChecksumsUnsupported = errors.New("version precedes checksums")
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Checksums                            */
/* -------------------------------------------------------------------------- */

// An interface for an 'Artifact' representing a checksums file.
type Checksums[T artifact.Artifact] interface {
	artifact.Artifact

	// Supports is a method to register which artifact the checksums are for.
	Supports(_ T)

	// Hash returns a 'hash.Hash' instance used to compute the checksum of the
	// supported artifact type.
	Hash() hash.Hash
}

/* -------------------------------------------------------------------------- */
/*                              Function: Compare                             */
/* -------------------------------------------------------------------------- */

// Compare takes a local executable archive and a local checksums file for
// executable archives and validates that the executable archive's checksum
// matches the expected value.
func Compare[T artifact.Artifact, U Checksums[T]](
	ctx context.Context,
	localArtifact artifact.Local[T],
	localChecksums artifact.Local[U],
) error {
	return checksumutil.Compare(
		ctx,
		localChecksums.Artifact.Hash(),
		localArtifact.Path,
		localChecksums.Path,
		localArtifact.Artifact.Name(),
	)
}
