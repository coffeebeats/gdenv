package archive

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/internal/progress"
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

	var reader io.Reader

	reader, err = xz.NewReader(f)
	if err != nil {
		return err
	}

	progressWriter, err := newTarProgressWriter(ctx, path)
	if err != nil {
		return err
	}

	if progressWriter != nil {
		reader = io.TeeReader(reader, progressWriter)
	}

	archive := tar.NewReader(reader)
	prefix := strings.TrimSuffix(filepath.Base(path), extensionTarXZ)

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
		out := filepath.Join(out, strings.TrimPrefix(hdr.Name, prefix+string(os.PathSeparator)))

		if err := extractFile(ctx, archive, hdr, out); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------- Function: extractFile ------------------------- */

// extractFile handles the extraction logic for each file in the Tar archive.
func extractFile(ctx context.Context, archive *tar.Reader, hdr *tar.Header, out string) error {
	mode := hdr.FileInfo().Mode()

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

	return nil
}

/* --------------------- Function: newTarProgressWriter --------------------- */

// newTarProgressWriter configures the 'progress.Progress' instance's
// 'total' bytes if one is found on the context.
func newTarProgressWriter(ctx context.Context, path string) (*progress.Writer, error) {
	p, ok := ctx.Value(progressKey{}).(*progress.Progress)
	if !ok || p == nil {
		return nil, nil //nolint:nilnil
	}

	sum, err := osutil.SizeOf(path)
	if err != nil {
		return nil, err
	}

	if err := p.Total(sum); err != nil {
		return nil, err
	}

	return progress.NewWriter(p), nil
}
