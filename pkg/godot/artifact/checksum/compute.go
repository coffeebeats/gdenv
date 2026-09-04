package checksum

import (
	"context"
	"hash"

	"github.com/coffeebeats/gdenv/internal/checksumutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Computes and returns the correct checksum of the specified archive.
func Compute[T artifact.Artifact](ctx context.Context, h hash.Hash, d artifact.Local[T]) (string, error) {
	return checksumutil.Compute(ctx, h, d.Path)
}
