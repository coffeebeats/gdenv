package checksum

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Source                               */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing a Godot source archive checksums file.
type Source checksums

// Compile-time verifications that 'Executable' implements 'Artifact'.
var _ artifact.Artifact = Source{} //nolint:exhaustruct

// Compile-time verifications that 'Source' implements 'Checksums'.
var _ Checksums[artifact.Folder[source.Source], source.Archive] = Source{} //nolint:exhaustruct

/* --------------------------- Function: NewSource -------------------------- */

// Returns a new 'Source' struct after validating the Godot version.
func NewSource(v version.Version) (Source, error) {
	var s Source

	if v.CompareNormal(versionSourceChecksumsSupported()) < 0 {
		return s, ErrChecksumsUnsupported
	}

	s.version = v

	return s, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (s Source) Name() string {
	return archive.TarXZ[source.Source]{Inner: source.New(s.version)}.Name() + ".sha256"
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (s Source) Version() version.Version {
	return s.version
}

/* ----------------------------- Impl: Checksums ---------------------------- */

func (s Source) supports(source.Archive) {} //nolint:unused

/* -------------------------------------------------------------------------- */
/*                  Function: versionSourceChecksumsSupported                 */
/* -------------------------------------------------------------------------- */

// Returns the first version at which 'godot-<version>.tar.xa.sha256' files for
// source archives began being published.
func versionSourceChecksumsSupported() version.Version {
	return version.MustParse("3.0")
}
