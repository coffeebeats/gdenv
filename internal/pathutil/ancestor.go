package pathutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

/* -------------------------------------------------------------------------- */
/*                             Function: Ancestor                             */
/* -------------------------------------------------------------------------- */

// Returns the closest ancestor of the specified 'path' which exists. If 'path'
// itself exists then it will be returned.
func Ancestor(path string) (string, error) {
	if path == "" {
		return "", fs.ErrInvalid
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		_, err := os.Stat(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return "", err
			}

			path = filepath.Dir(path)
			if path == "" {
				return "", fs.ErrNotExist
			}

			continue
		}

		return path, nil
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: AncestorDir                           */
/* -------------------------------------------------------------------------- */

// Returns the closest ancestor directory of the specified 'path' which exists.
// If 'path' itself exists and is a directory then it will be returned.
func AncestorDir(path string) (string, error) {
	if path == "" {
		return "", fs.ErrInvalid
	}

	path, err := Ancestor(path)
	if err != nil {
		return "", err
	}

	return filepath.Dir(path), nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: AncestorMode                           */
/* -------------------------------------------------------------------------- */

// Returns the 'fs.FileMode' of the closest ancestor directory of the specified
// 'path' which exists. If 'path' itself exists and is a directory then it will
// be returned.
func AncestorMode(path string) (fs.FileMode, error) {
	ancestor, err := AncestorDir(path)
	if err != nil {
		return 0, err
	}

	info, err := os.Stat(ancestor)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}
