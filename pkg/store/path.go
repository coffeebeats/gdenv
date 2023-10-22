package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const envStore = "GDENV_HOME"
const storeName = "gdenv"

var (
	ErrIllegalPath   = errors.New("illegal store path")
	ErrInvalidPath   = errors.New("invalid file path")
	ErrMissingEnvVar = errors.New("missing environment variable")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Path                               */
/* -------------------------------------------------------------------------- */

// Returns the user-configured path to the 'gdenv' store.
func Path() (string, error) {
	path := os.Getenv(envStore)
	if path == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingEnvVar, envStore)
	}

	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("%w; expected absolute path: %s", ErrInvalidPath, path)
	}

	if base := filepath.Base(path); base != storeName && base != "."+storeName {
		return "", fmt.Errorf("%w: '%s'", ErrIllegalPath, path)
	}

	return filepath.Clean(path), nil
}
