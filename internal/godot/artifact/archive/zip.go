package archive

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
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

	// Extract all files within the archive.
	for _, f := range archive.File {
		mode := f.FileInfo().Mode()
		out := filepath.Join(out, f.Name) //nolint:gosec

		if f.FileInfo().IsDir() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err := os.MkdirAll(out, mode); err != nil {
				return err
			}

			continue
		}

		src, err := archive.Open(f.Name)
		if err != nil {
			return err
		}

		if err := copyFile(ctx, src, mode, out); err != nil {
			return err
		}
	}

	return nil
}
