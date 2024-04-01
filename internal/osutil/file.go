package osutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/ioutil"
)

// Only write to 'out'; create a new file/overwrite an existing.
const copyFileWriteFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

/* -------------------------------------------------------------------------- */
/*                             Function: CopyFile                             */
/* -------------------------------------------------------------------------- */

// CopyFile is a utility function for copying an 'io.Reader' to a new file
// created with the specified 'os.FileMode'.
func CopyFile(ctx context.Context, src, out string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(out, copyFileWriteFlag, info.Mode())
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
		return err
	}

	return nil
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

/* -------------------------------------------------------------------------- */
/*                            Function: ForceRename                           */
/* -------------------------------------------------------------------------- */

// Moves a file or directory from one path to a new path.
//
// NOTE: This will overwrite any file or directory which already exists at
// the 'newPath' parameter.
func ForceRename(oldPath, newPath string) error {
	// Verify that the file-to-add exists.
	if _, err := os.Stat(oldPath); err != nil {
		return err
	}

	// Verify that the new path either is a file or doesn't exist.
	info, err := os.Stat(newPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	// If the target is already a directory, remove it first.
	if info != nil && info.IsDir() {
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}
	}

	return os.Rename(oldPath, newPath)
}

/* -------------------------------------------------------------------------- */
/*                              Function: ModeOf                              */
/* -------------------------------------------------------------------------- */

// ModeOf returns the file mode of the specified file.
func ModeOf(path string) (fs.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: SizeOf                              */
/* -------------------------------------------------------------------------- */

// SizeOf returns the size of the specified file in bytes.
func SizeOf(path string) (uint64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return uint64(info.Size()), nil
}
