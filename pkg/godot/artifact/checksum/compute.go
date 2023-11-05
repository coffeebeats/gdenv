package checksum

import (
	"context"
	"encoding/hex"
	"hash"
	"io"
	"os"

	"github.com/coffeebeats/gdenv/internal/ioutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Computes and returns the correct checksum of the specified archive.
func Compute[T artifact.Artifact](ctx context.Context, h hash.Hash, d artifact.Local[T]) (string, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(h, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
