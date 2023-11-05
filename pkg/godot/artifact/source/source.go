package source

import (
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	namePrefix    = "godot"
	nameSeparator = "-"
)

type Archive = archive.TarXZ[Source]

/* -------------------------------------------------------------------------- */
/*                               Struct: Source                               */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing Godot source code for a specific version.
type Source struct {
	version version.Version
}

// Compile-time verifications that 'Source' implements 'Artifact'.
var _ artifact.Artifact = (*Source)(nil)

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Source' for the specified 'Version'.
func New(v version.Version) Source {
	return Source{v}
}

/* ------------------------ Impl: archive.Archivable ------------------------ */

// Allows 'Source' to be used by 'Archive' implementation.
func (s Source) Archivable() {}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'Source' as a Godot release artifact.
func (s Source) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

// Returns the name of the Godot source directory for the specified 'Version'.
//
// NOTE: Godot names its executables in the format 'godot-<VERSION>'.
func (s Source) Name() string {
	var name strings.Builder

	name.WriteString(namePrefix)
	name.WriteString(nameSeparator)

	name.WriteString(strings.TrimPrefix(s.version.String(), version.Prefix))

	return name.String()
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (s Source) Version() version.Version {
	return s.version
}

/* ----------------------- Impl: checksum.Checksumable ---------------------- */

// Checksumable "registers" 'Source' as a Godot release artifact with published
// file checksums.
func (s Source) Checksumable() {}

/* ----------------------------- Impl: Stringer ----------------------------- */

func (s Source) String() string {
	return s.Name()
}
