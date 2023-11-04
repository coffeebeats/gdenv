package checksum

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/coffeebeats/gdenv/internal/ioutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
)

var ErrUnsupportedArtifact = errors.New("unsupported artifact")

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Computes and returns the correct checksum of the specified archive.
func Compute[T archive.Archive](ctx context.Context, d artifact.Local[T]) (string, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	var h hash.Hash

	switch any(d.Artifact).(type) { // FIXME: https://github.com/golang/go/issues/45380
	case executable.Archive:
		h = sha512.New()
		if _, err := io.Copy(h, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
			return "", err
		}

	case source.Archive:
		h = sha256.New()
		if _, err := io.Copy(h, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("%w: '%T'", ErrUnsupportedArtifact, d.Artifact)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}