package osutil

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrUnsupportedFileType is returned when attempting to copy a file type that's
// not supported.
var ErrUnsupportedFileType = errors.New("unsupported file type")

/* -------------------------------------------------------------------------- */
/*                              Function: CopyDir                             */
/* -------------------------------------------------------------------------- */

// CopyDir recursively copies a directory from 'srcDir' to 'dstDir', preserving
// soft links. All regular files will be hard copied. Note that file attributes
// are not preserved, so this should only be used when the folder contents are
// required in the original structure. This implementation is based on [1].
//
// [1] https://github.com/moby/moby/blob/master/daemon/graphdriver/copy/copy.go
func CopyDir(ctx context.Context, srcDir, dstDir string) error { //nolint:cyclop
	return filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		switch mode := f.Mode(); {
		case mode.IsRegular():
			if err := os.Link(srcPath, dstPath); err != nil {
				return err
			}

		case mode.IsDir():
			if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
				return err
			}

		case mode&os.ModeSymlink != 0:
			link, err := os.Readlink(srcPath)
			if err != nil {
				return err
			}

			if err := os.Symlink(link, dstPath); err != nil {
				return err
			}

		default:
			return fmt.Errorf("%w: (%d / %s) for %s", ErrUnsupportedFileType, f.Mode(), f.Mode().String(), srcPath)
		}

		return nil
	})
}

/* -------------------------------------------------------------------------- */
/*                             Function: EnsureDir                            */
/* -------------------------------------------------------------------------- */

// EnsureDir verifies that the specified path exists, is a directory, and has
// the specified permission bits set.
func EnsureDir(path string, perm fs.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err := os.MkdirAll(path, perm); err != nil {
			return err
		}
	}

	if info != nil {
		if !info.IsDir() {
			return fmt.Errorf("%w: %s", fs.ErrExist, path)
		}

		if info.Mode().Perm()&perm == 0 {
			return os.Chmod(path, info.Mode()|perm)
		}
	}

	return nil
}
