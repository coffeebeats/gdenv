package store

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrEnvVariable = errors.New("store: environment variable 'GDENV_HOME' not defined")
	ErrMissingPath = errors.New("store: missing file path")
	ErrInvalidPath = errors.New("store: invalid file path")

	envVarStore = "GDENV_HOME"
)

/* ----------------------------- Function: Clean ---------------------------- */

// Returns a "cleaned" version of the specified store path.
func Clean(p string) (string, error) {
	if p == "" {
		return p, ErrMissingPath
	}

	p, err := filepath.Abs(p)
	if err != nil {
		return p, errors.Join(ErrInvalidPath, err)
	}

	return p, nil
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
		return "", ErrEnvVariable
	}

	return filepath.Clean(p), nil
}
