package source

import (
	"errors"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	namePrefix    = "godot"
	nameSeparator = "-"
)

var ErrInvalidPlatform = errors.New("invalid platform")

/* -------------------------------------------------------------------------- */
/*                               Struct: Source                               */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing Godot source code for a specific version.
type Source struct {
	version version.Version
}

// Compile-time verifications that 'Source' implements 'Artifact'.
var _ artifact.Artifact = Source{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Source' for the specified 'Version'.
func New(v version.Version) Source {
	return Source{v}
}

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
