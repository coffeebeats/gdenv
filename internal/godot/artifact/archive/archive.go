package archive

import "github.com/coffeebeats/gdenv/internal/godot/artifact"

// An interface representing a compressed 'Artifact' archive.
type Archive[T artifact.Artifact] interface {
	artifact.Artifact

	Contents() T
	extract(path, out string) (artifact.Local[T], error)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Given a downloaded 'Archive', extract the contents and return a local
// 'Artifact' pointing to it.
func Extract[T artifact.Artifact](a artifact.Local[Archive[T]], out string) (artifact.Local[T], error) {
	return a.Artifact.extract(a.Path, out)
}
