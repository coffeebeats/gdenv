package archive

import (
	"archive/zip"
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/internal/progress"
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
// NOTE: This method does not detect insecure filepaths included in the archive.
// Instead, ensure the binary is compiled with the GODEBUG option
// 'zipinsecurepath=0' (see https://github.com/golang/go/issues/55356).
func (a Zip[T]) extract(ctx context.Context, path, out string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer archive.Close()

	baseDirMode, err := osutil.ModeOf(out)
	if err != nil {
		return err
	}

	// There doesn't appear to be a good way to read the compressed bytes during
	// extraction. Instead, use a manual writer and record progress in steps
	// after each file completes.
	progressWriter, err := newZipProgressWriter(ctx, path)
	if err != nil {
		return err
	}

	// Extract all files within the archive.
	for _, f := range archive.File {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := extractZipFile(ctx, archive, f, out, baseDirMode); err != nil {
			return err
		}

		if progressWriter != nil {
			progressWriter.Add(f.CompressedSize64)
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
	out = filepath.Join(out, f.Name) //nolint:gosec

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

/* --------------------- Function: newZipProgressWriter --------------------- */

// newZipProgressWriter configures the 'progress.Progress' instance's 'total'
// bytes if one is found on the context. If one is found then a valid pointer to
// a 'progress.ManualWriter' will be returned.
//
// NOTE: Using a pointer for optionality here is not ideal, but there isn't much
// benefit to improving this.
func newZipProgressWriter(ctx context.Context, path string) (*progress.ManualWriter, error) {
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

	if err := p.Total(sum); err != nil {
		return nil, err
	}

	return progress.NewManualWriter(p), nil
}
