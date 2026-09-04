package archive

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/extract"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const extensionTarXZ = ".tar.xz"

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

// A struct representing an 'XZ'-compressed archive.
type TarXZ[T Archivable] struct {
	Inner T
}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'TarXZ' as a Godot release artifact.
func (a TarXZ[T]) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

func (a TarXZ[T]) Name() string {
	name := a.Inner.Name()
	if name != "" {
		name += extensionTarXZ
	}

	return name
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (a TarXZ[T]) Version() version.Version {
	return a.Inner.Version()
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
//
// NOTE: While this method *does* detect insecure filepaths included in the
// archive using the same method implemented by Go, this binary should still be
// compiled with the GODEBUG option 'tarinsecurepath=0' in the event that the
// implementation changes (see https://github.com/golang/go/issues/55356).
func (a TarXZ[T]) extract(ctx context.Context, path, out string) error {
	// Remove the name of the tar-file from each archived filepath; this is to
	// facilitate extracting contents directly into the 'out' path.
	prefix := strings.TrimSuffix(filepath.Base(path), extensionTarXZ)

	return extract.TarXZ(ctx, path, out, extract.WithStripPrefix(prefix))
}
