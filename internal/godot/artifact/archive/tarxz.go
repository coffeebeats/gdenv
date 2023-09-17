package archive

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
)

const extensionTarXZ = ".tar.xz"

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

// A struct representing an 'XZ'-compressed archive.
type TarXZ[T artifact.Artifact] struct {
	Inner T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a TarXZ[T]) Name() string {
	return a.Inner.Name() + extensionTarXZ
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Returns the 'Artifact' contained in this 'Archive'.
func (a TarXZ[T]) Contents() T { //nolint:ireturn
	return a.Inner
}

// Extracts the archive to the specified file path.
func (a TarXZ[T]) extract(path, out string) (artifact.Local[T], error) {
	// TODO: Implement the archive extraction.
	return artifact.Local[T]{Artifact: a.Inner, Path: out}, nil
}
