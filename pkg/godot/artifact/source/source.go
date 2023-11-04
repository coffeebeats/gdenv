package source

import (
	"errors"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	namePrefix    = "godot"
	nameSeparator = "-"
)

var ErrInvalidPlatform = errors.New("invalid platform")

/* -------------------------------------------------------------------------- */
/*                               Struct: Source                               */
/* -------------------------------------------------------------------------- */

type Archive = archive.TarXZ[Source]

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

/* ---------------------------- Impl: Archivable ---------------------------- */

// Allows 'Source' to be used by 'Archive' implementation.
func (s Source) Archivable() {}

/* ----------------------------- Impl: Artifact ----------------------------- */

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

/* ----------------------------- Impl: Versioned ---------------------------- */

func (s Source) Version() version.Version {
	return s.version
}

/* ----------------------------- Impl: Stringer ----------------------------- */

func (s Source) String() string {
	return s.Name()
}
