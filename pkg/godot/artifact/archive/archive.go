package archive

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/internal/extract"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

// ErrExtractFailed is returned when an archive's contents cannot be extracted.
//
// NOTE: This is an alias of the error returned by the underlying extraction
// implementation so that 'errors.Is' checks continue to match.
var ErrExtractFailed = extract.ErrFailed

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithProgress creates a sub-context with an associated progress reporter. The
// result can be passed to the extract function(s) in this package to get
// updates on extraction progress.
func WithProgress(ctx context.Context, p *progress.Progress) context.Context {
	return extract.WithProgress(ctx, p)
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

	// Validate that the 'out' parameter exists and is a directory.
	info, err := os.Stat(out)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: expected a directory", fs.ErrInvalid)
	}

	log.Debugf("extracting archive %s: %s", filepath.Base(a.Path), out)

	// Extract the contents to the specified 'out' directory.
	return a.Artifact.extract(ctx, a.Path, out)
}
