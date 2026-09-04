package archive

import (
	"context"

	"github.com/coffeebeats/gdenv/internal/extract"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip'-compressed archive.
type Zip[T Archivable] struct {
	Inner T
}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'Zip' as a Godot release artifact.
func (a Zip[T]) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

func (a Zip[T]) Name() string {
	name := a.Inner.Name()
	if name != "" {
		name += extensionZip
	}

	return name
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (a Zip[T]) Version() version.Version {
	return a.Inner.Version()
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: While this method *does* detect insecure filepaths included in the
// archive using the same method implemented by Go, this binary should still be
// compiled with the GODEBUG option 'zipinsecurepath=0' in the event that the
// implementation changes (see https://github.com/golang/go/issues/55356).
func (a Zip[T]) extract(ctx context.Context, path, out string) error {
	return extract.Zip(ctx, path, out)
}
