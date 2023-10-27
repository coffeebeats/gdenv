package archive

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

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
func (a TarXZ[T]) extract(ctx context.Context, path, out string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	prefix := strings.TrimSuffix(filepath.Base(path), extensionTarXZ)

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

		// Remove the name of the tar-file from the filepath; this is to
		// facilitate extracting contents directly into the 'out' path.
		name := strings.TrimPrefix(hdr.Name, prefix+string(os.PathSeparator))

		mode := hdr.FileInfo().Mode()
		out := filepath.Join(out, name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if ctx.Err() != nil {
				return ctx.Err()
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
