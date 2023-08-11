package pin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const pinFilename = ".godot-version"

var (
	ErrMissingPath = errors.New("pin: missing file path")
	ErrInvalidPath = errors.New("pin: invalid file path")
)

/* ----------------------------- Function: Clean ---------------------------- */

// Returns a "cleaned" version of the specified pin file path.
func Clean(p string) (string, error) {
	if p == "" {
		return p, ErrMissingPath
	}

	p, err := filepath.Abs(p)
	if err != nil {
		return p, errors.Join(ErrInvalidPath, err)
	}

	fmt.Println(p)

	if filepath.Base(p) != pinFilename {
		p = filepath.Join(p, pinFilename)
	}

	return p, nil
}

/* ---------------------------- Function: Exists ---------------------------- */

// Returns whether the specified pin file exists.
func Exists(p string) bool {
	p, err := Clean(p)
	if err != nil {
		return false
	}

	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
