package archive

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

const extensionTarXZ = ".tar.xz"

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

// A struct representing an 'XZ'-compressed archive.
type TarXZ[T Archivable] struct {
	Inner T
}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'TarXZ' as a Godot release artifact.
func (a TarXZ[T]) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

func (a TarXZ[T]) Name() string {
	name := a.Inner.Name()
	if name != "" {
		name += extensionTarXZ
	}

	return name
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (a TarXZ[T]) Version() version.Version {
	return a.Inner.Version()
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: While this method *does* detect insecure filepaths included in the
// archive using the same method implemented by Go, this binary should still be
// compiled with the GODEBUG option 'tarinsecurepath=0' in the event that the
// implementation changes (see https://github.com/golang/go/issues/55356).
func (a TarXZ[T]) extract(ctx context.Context, path, out string) error { //nolint:cyclop,funlen
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	reader, err := newFileReaderWithProgress(ctx, f)
	if err != nil {
		return err
	}

	reader, err = xz.NewReader(reader)
	if err != nil {
		return err
	}

	archive := tar.NewReader(reader)

	baseDirMode, err := osutil.ModeOf(out)
	if err != nil {
		return err
	}

	prefix := strings.TrimSuffix(filepath.Base(path), extensionTarXZ)

	// Extract all files within the archive.
	for {
		hdr, err := archive.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}

			closeProgress(ctx)

			break
		}

		name := hdr.Name

		// See https://cs.opensource.google/go/go/+/refs/tags/go1.21.3:src/archive/tar/reader.go;l=60-67.
		if !filepath.IsLocal(name) || strings.Contains(name, `\`) || strings.Contains(name, "..") {
			return fmt.Errorf("%w: %s", tar.ErrInsecurePath, name)
		}

		// Remove the name of the tar-file from the filepath; this is to
		// facilitate extracting contents directly into the 'out' path.
		name = strings.TrimPrefix(name, prefix+"/") // Archive always uses the '/' separator.
		if strings.HasPrefix(name, prefix) {
			return fmt.Errorf(
				"%w: couldn't trim prefix: %s from %s",
				ErrExtractFailed,
				prefix, name,
			)
		}

		if err := extractTarFile(ctx, archive, hdr, out, baseDirMode); err != nil {
			return err
		}
	}

	return nil
}

/* ------------------------- Function: closeProgress ------------------------ */

// closeProgress updates the 'progress.Progress' instance attached to the
// context to 100% complete. This is because the 'tar.Reader' can discard bytes
// from the last file, causing the reported progress to not be accurate at
// close. There doesn't seem to be a way to get the exact amount, so just add
// what's missing.
func closeProgress(ctx context.Context) {
	p, ok := ctx.Value(progressKey{}).(*progress.Progress)
	if ok && p != nil {
		if remaining := p.Total() - p.Current(); remaining > 0 {
			p.Add(remaining)
		}
	}
}

/* ------------------------ Function: extractTarFile ------------------------ */

// extractFile handles the extraction logic for each file in the Tar archive.
func extractTarFile(
	ctx context.Context,
	archive *tar.Reader,
	hdr *tar.Header,
	out string,
	baseDirMode fs.FileMode,
) error {
	// Ensure the parent directory exists with best-effort permissions. If
	// the zip archive already contains the directory as an entry then this
	// will have no effect.
	if err := os.MkdirAll(filepath.Dir(out), baseDirMode); err != nil {
		return err
	}

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
