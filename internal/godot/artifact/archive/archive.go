package archive

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Archive                             */
/* -------------------------------------------------------------------------- */

// An alias for a locally-available 'Archive'.
type Local = artifact.Local[Archive]

// An interface representing a compressed 'Artifact' archive.
type Archive interface {
	artifact.Artifact

	extract(path, out string) error
}

/* -------------------------------------------------------------------------- */
/*                            Interface: Archivable                           */
/* -------------------------------------------------------------------------- */

// An interface representing an 'Artifact' that can be compressed into an
// archive.
type Archivable interface {
	artifact.Artifact

	Archivable()
}

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Given a downloaded 'Archive', extract the contents and return a local
// 'Artifact' pointing to it.
func Extract[T Archive](a artifact.Local[T], out string) error {
	// Validate that the artifact exists.
	ok, err := a.Exists()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("%w: '%s'", fs.ErrNotExist, a.Path)
	}

	// Validate that the 'out' parameter either doesn't exist or is a directory.
	info, err := os.Stat(out)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("%w: expected a directory", fs.ErrInvalid)
	}

	// Extract the contents to the specified 'out' directory.
	return a.Artifact.extract(a.Path, out)
}
