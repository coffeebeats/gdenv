package checksum

import (
	"context"

	"github.com/coffeebeats/gdenv/internal/checksumutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

// NOTE: These are aliases of the errors returned by the underlying checksum
// implementation so that 'errors.Is' checks continue to match.
var (
	ErrChecksumNotFound    = checksumutil.ErrNotFound
	ErrConflictingChecksum = checksumutil.ErrConflicting
	ErrUnrecognizedFormat  = checksumutil.ErrUnrecognizedFormat
)

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Given a locally-available checksums file, find and return the checksum for
// the specified archive.
func Extract[T artifact.Artifact, U Checksums[T]](
	ctx context.Context,
	local artifact.Local[U],
	a T,
) (string, error) {
	return checksumutil.Lookup(ctx, local.Path, a.Name())
}
