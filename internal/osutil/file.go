package osutil

import (
	"errors"
	"io/fs"
	"os"
)

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
