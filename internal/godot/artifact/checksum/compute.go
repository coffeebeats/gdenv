package checksum

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
)

var (
	ErrUnsupportedArtifact = errors.New("unsupported artifact")
)

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Computes and returns the sha-512 checksum of the specified archive.
func Compute[T artifact.Artifact, U archive.Archive[T]](d artifact.Local[U]) (string, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	var h hash.Hash

	contents := d.Artifact.Contents()

	switch any(contents).(type) { // FIXME: https://github.com/golang/go/issues/45380
	case executable.Executable:
		h = sha512.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}

	case source.Source:
		h = sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("%w: '%T'", ErrUnsupportedArtifact, contents)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
