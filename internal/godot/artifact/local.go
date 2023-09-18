package artifact

import (
	"errors"
	"os"
)

var ErrMissingPath = errors.New("missing path")

type Downloaded = Local[Artifact]

/* -------------------------------------------------------------------------- */
/*                                Struct: Local                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is locally-available on the file system.
type Local[T Artifact] struct {
	Artifact T
	Path     string
}

/* ----------------------------- Method: Exists ----------------------------- */

// Returns whether the downloaded file exists on the local file system.
func (l Local[T]) Exists() (bool, error) {
	if l.Path == "" {
		return false, ErrMissingPath
	}

	if _, err := os.Stat(l.Path); err != nil {
		return false, err
	}

	return true, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (l Local[T]) Name() string {
	return l.Artifact.Name()
}
