package archive

import (
	"archive/zip"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip'-compressed archive.
type Zip[T Archivable] struct {
	Artifact T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Zip[T]) Name() string {
	name := a.Artifact.Name()
	if name != "" {
		name += extensionZip
	}

	return name
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: While this method *does* detect insecure filepaths included in the
// archive using the same method implemented by Go, this binary should still be
// compiled with the GODEBUG option 'zipinsecurepath=0' in the event that the
// implementation changes (see https://github.com/golang/go/issues/55356).
func (a Zip[T]) extract(ctx context.Context, path, out string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer archive.Close()

	// There doesn't appear to be a good way to read the compressed bytes during
	// extraction. Instead, use a manual writer and record progress in steps
	// after each file completes.
	p, err := newZipProgress(ctx, path)
	if err != nil {
		return err
	}

	baseDirMode, err := osutil.ModeOf(out)
	if err != nil {
		return err
	}

	// Extract all files within the archive.
	for _, f := range archive.File {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		name := f.Name

		// See https://cs.opensource.google/go/go/+/refs/tags/go1.21.3:src/archive/zip/reader.go;l=168-173.
		if !filepath.IsLocal(name) || strings.Contains(name, `\`) || strings.Contains(name, "..") {
			return fmt.Errorf("%w: %s", zip.ErrInsecurePath, name)
		}

		out := filepath.Join(out, name) //nolint:gosec

		if err := extractZipFile(ctx, archive, f, out, baseDirMode); err != nil {
			return err
		}

		if p != nil {
			p.Add(f.CompressedSize64)
		}
	}

	return nil
}

/* ------------------------ Function: extractZipFile ------------------------ */

// extractZipFile extracts a single entry from a zip archive to the specified
// destination path.
func extractZipFile(
	ctx context.Context,
	a *zip.ReadCloser,
	f *zip.File,
	out string,
	baseDirMode fs.FileMode,
) error {
	// Ensure the parent directory exists with best-effort permissions. If
	// the zip archive already contains the directory as an entry then this
	// will have no effect.
	if err := os.MkdirAll(filepath.Dir(out), baseDirMode); err != nil {
		return err
	}

	mode := f.FileInfo().Mode()

	// Create all the ancestor directories if required.
	if err := os.MkdirAll(filepath.Dir(out), mode); err != nil {
		return err
	}

	if f.FileInfo().IsDir() {
		if err := os.Mkdir(out, mode); err != nil {
			return err
		}

		return nil
	}

	src, err := a.Open(f.Name)
	if err != nil {
		return err
	}

	if err := copyFile(ctx, src, mode, out); err != nil {
		return err
	}

	return nil
}

/* ------------------------ Function: newZipProgress ------------------------ */

// newZipProgress sets the 'total' value of the 'progress.Progress' instance
// attached to the context, if one exists. A pointer to the provided
// 'progress.Progress' is returned.
//
// NOTE: Using a pointer for optionality here is not ideal, but there isn't much
// benefit to improving this.
func newZipProgress(ctx context.Context, path string) (*progress.Progress, error) {
	p, ok := ctx.Value(progressKey{}).(*progress.Progress)
	if !ok || p == nil {
		return nil, nil //nolint:nilnil
	}

	var sum uint64

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	for _, f := range r.File {
		sum += f.CompressedSize64
	}

	if err := p.SetTotal(sum); err != nil {
		return nil, err
	}

	return p, nil
}
