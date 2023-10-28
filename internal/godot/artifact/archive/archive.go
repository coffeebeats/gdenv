package archive

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/ioutil"
	"github.com/coffeebeats/gdenv/internal/progress"
)

// Only write to 'out'; create a new file/overwrite an existing.
const copyFileWriteFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

type progressKey struct{}

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithProgress creates a sub-context with an associated progress reporter. The
// result can be passed to the extract function(s) in this package to get
// updates on extraction progress.
func WithProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKey{}, p)
}

/* -------------------------------------------------------------------------- */
/*                             Interface: Archive                             */
/* -------------------------------------------------------------------------- */

// An alias for a locally-available 'Archive'.
type Local = artifact.Local[Archive]

// An interface representing a compressed 'Artifact' archive.
type Archive interface {
	artifact.Artifact

	extract(ctx context.Context, path, out string) error
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
func Extract[T Archive](ctx context.Context, a artifact.Local[T], out string) error {
	// Validate that the artifact exists.
	if !a.Exists() {
		return fmt.Errorf("%w: '%s'", fs.ErrNotExist, a.Path)
	}

	// Validate that the 'out' parameter either doesn't exist or is a directory.
	info, err := os.Stat(out)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: expected a directory", fs.ErrInvalid)
	}

	// Extract the contents to the specified 'out' directory.
	return a.Artifact.extract(ctx, a.Path, out)
}

/* --------------------------- Function: copyFile --------------------------- */

// A shared helper function which copies the contents of an 'io.Reader' to a new
// file created with the specified 'os.FileMode'.
func copyFile(ctx context.Context, f io.Reader, mode fs.FileMode, out string) error {
	dst, err := os.OpenFile(out, copyFileWriteFlag, mode)
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, ioutil.NewReaderClosure(ctx, f.Read)); err != nil {
		return err
	}

	return nil
}
