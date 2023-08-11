package pin

import (
	"errors"
	"os"
	"path/filepath"
)

const pinFilename = ".godot-version"

var (
	ErrInvalidPath = errors.New("pin: invalid file path")
	ErrMissingPath = errors.New("pin: missing file path")
)

/* ----------------------------- Function: Clean ---------------------------- */

// Returns a "cleaned" version of the specified pin file path.
func Clean(path string) (string, error) {
	if path == "" {
		return path, ErrMissingPath
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return path, errors.Join(ErrInvalidPath, err)
	}

	if filepath.Base(path) != pinFilename {
		path = filepath.Join(path, pinFilename)
	}

	return path, nil
}

/* ---------------------------- Function: Exists ---------------------------- */

// Returns whether the specified pin file exists.
func Exists(path string) bool {
	path, err := Clean(path)
	if err != nil {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
