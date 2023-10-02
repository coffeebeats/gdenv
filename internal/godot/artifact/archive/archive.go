package archive

import "github.com/coffeebeats/gdenv/internal/godot/artifact"

/* -------------------------------------------------------------------------- */
/*                             Interface: Archive                             */
/* -------------------------------------------------------------------------- */

// An interface representing a compressed 'Artifact' archive.
type Archive interface {
	artifact.Artifact

	extract(path, out string) error
}

/* -------------------------------------------------------------------------- */
/*                            Interface: Archivable                           */
/* -------------------------------------------------------------------------- */

// An interface representing an 'Artifact' that can be archived.
type Archivable interface {
	artifact.Artifact

	Archivable()
}

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Given a downloaded 'Archive', extract the contents and return a local
// 'Artifact' pointing to it.
func Extract[T Archive](a artifact.Local[Archive], out string) error {
	return a.Artifact.extract(a.Path, out)
}
