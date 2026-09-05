package extract

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

// A decompressor wraps a compressed input stream in a reader which yields the
// decompressed 'tar' stream.
type decompressor func(io.Reader) (io.Reader, error)

/* -------------------------------------------------------------------------- */
/*                               Function: TarXZ                              */
/* -------------------------------------------------------------------------- */

// TarXZ extracts the contents of an 'XZ'-compressed tarball into the specified
// directory.
func TarXZ(ctx context.Context, path, out string, opts ...Option) error {
	return extractTar(ctx, path, out, newXZReader, opts...)
}

/* -------------------------------------------------------------------------- */
/*                               Function: TarGz                              */
/* -------------------------------------------------------------------------- */

// TarGz extracts the contents of a 'gzip'-compressed tarball into the specified
// directory.
func TarGz(ctx context.Context, path, out string, opts ...Option) error {
	return extractTar(ctx, path, out, newGzipReader, opts...)
}

/* --------------------------- Function: newXZReader ------------------------ */

func newXZReader(r io.Reader) (io.Reader, error) {
	return xz.NewReader(r)
}

/* -------------------------- Function: newGzipReader ----------------------- */

func newGzipReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

/* -------------------------------------------------------------------------- */
/*                            Function: extractTar                            */
/* -------------------------------------------------------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: While this function *does* detect insecure filepaths included in the
// archive using the same method implemented by Go, this binary should still be
// compiled with the GODEBUG option 'tarinsecurepath=0' in the event that the
// implementation changes (see https://github.com/golang/go/issues/55356).
func extractTar( //nolint:cyclop
	ctx context.Context,
	path, out string,
	decompress decompressor,
	opts ...Option,
) error {
	cfg := newConfig(opts...)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	reader, err := newFileReaderWithProgress(ctx, f)
	if err != nil {
		return err
	}

	reader, err = decompress(reader)
	if err != nil {
		return err
	}

	archive := tar.NewReader(reader)

	baseDirMode, err := osutil.ModeOf(out)
	if err != nil {
		return err
	}

	// Extract all files within the archive.
	for {
		hdr, err := archive.Next()
		if err != nil {
			if err != io.EOF { //nolint:errorlint
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

		name, err = stripPrefix(name, cfg.stripPrefix)
		if err != nil {
			return err
		}

		out := filepath.Join(out, name) //nolint:gosec

		if err := extractTarFile(ctx, archive, hdr, out, baseDirMode); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------------- Function: stripPrefix ------------------------ */

// stripPrefix removes the name of the enclosing directory from the archived
// filepath; this is to facilitate extracting contents directly into the output
// path. An empty prefix leaves the name unchanged.
func stripPrefix(name, prefix string) (string, error) {
	if prefix == "" {
		return name, nil
	}

	name = strings.TrimPrefix(name, prefix+"/") // Archive always uses the '/' separator.
	if strings.HasPrefix(name, prefix) {
		return "", fmt.Errorf(
			"%w: couldn't trim prefix: %s from %s",
			ErrFailed,
			prefix, name,
		)
	}

	return name, nil
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

// extractTarFile handles the extraction logic for each file in the Tar archive.
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
