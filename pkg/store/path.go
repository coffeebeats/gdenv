package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const envVarStore = "GDENV_HOME"

var (
	ErrInvalidPath   = errors.New("store: invalid file path")
	ErrMissingEnvVar = errors.New(fmt.Sprintf("store: environment variable '%s' not defined", envVarStore))
	ErrMissingPath   = errors.New("store: missing file path")
)

/* ----------------------------- Function: Clean ---------------------------- */

// Returns a "cleaned" version of the specified store path.
func Clean(path string) (string, error) {
	if path == "" {
		return "", ErrMissingPath
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Join(ErrInvalidPath, err)
	}

	return path, nil
}

/* ---------------------------- Function: Exists ---------------------------- */

// Returns whether the specified store path exists.
func Exists(p string) bool {
	p, err := Clean(p)
	if err != nil {
		return false
	}

	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	return info.IsDir()
}

/* ----------------------------- Function: Path ----------------------------- */

// Returns the user-configured path to the 'gdenv' store.
func Path() (string, error) {
	p := os.Getenv(envVarStore)
	if p == "" {
		return "", ErrMissingEnvVar
	}

	if !filepath.IsAbs(p) {
		return "", fmt.Errorf("%w; expected absolute path: %s", ErrInvalidPath, p)
	}

	return filepath.Clean(p), nil
}
