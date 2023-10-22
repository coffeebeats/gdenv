package pin

import (
	"errors"
	"path/filepath"
)

const pinFilename = ".godot-version"

var (
	ErrInvalidPath = errors.New("invalid file path")
	ErrMissingPath = errors.New("missing file path")
)

/* -------------------------------------------------------------------------- */
/*                               Function: clean                              */
/* -------------------------------------------------------------------------- */

// Returns a "cleaned" version of the specified pin file path.
func clean(path string) (string, error) {
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
