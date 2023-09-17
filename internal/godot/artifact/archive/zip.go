package archive

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
)

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip'-compressed archive.
type Zip[T artifact.Artifact] struct {
	Inner T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Zip[T]) Name() string {
	return a.Inner.Name() + extensionZip
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Returns the 'Artifact' contained in this 'Archive'.
func (a Zip[T]) Contents() T { //nolint:ireturn
	return a.Inner
}

// Extracts the archive to the specified file path.
func (a Zip[T]) extract(path, out string) (artifact.Local[T], error) {
	// TODO: Implement the archive extraction.
	return artifact.Local[T]{Artifact: a.Inner, Path: out}, nil
}
