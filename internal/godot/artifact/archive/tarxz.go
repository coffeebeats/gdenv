package archive

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

const extensionTarXZ = ".tar.xz"

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

// A struct representing an 'XZ'-compressed archive.
type TarXZ[T Archivable] struct {
	Artifact T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a TarXZ[T]) Name() string {
	name := a.Artifact.Name()
	if name != "" {
		name += extensionTarXZ
	}

	return name
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: This method does not detect insecure filepaths included in the archive.
// Instead, ensure the binary is compiled with the GODEBUG option
// 'tarinsecurepath=0' (see https://github.com/golang/go/issues/55356).
func (a TarXZ[T]) extract(ctx context.Context, path, out string) error { //nolint:cyclop
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	decompressed, err := xz.NewReader(f)
	if err != nil {
		return err
	}

	archive := tar.NewReader(decompressed)

	// Extract all files within the archive.
	for {
		hdr, err := archive.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		mode := hdr.FileInfo().Mode()
		out := filepath.Join(out, hdr.Name) //nolint:gosec

		switch hdr.Typeflag {
		case tar.TypeDir:
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err := os.MkdirAll(out, mode); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := copyFile(ctx, archive, mode, out); err != nil {
				return err
			}
		}
	}

	return nil
}
