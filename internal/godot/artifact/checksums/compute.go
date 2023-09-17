package checksums

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Computes and returns the sha-512 checksum of the specified file.
func Compute(d artifact.Local[archive.Archive]) (string, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
