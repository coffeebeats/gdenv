package executable

import (
	"errors"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	namePrefix    = "Godot"
	nameSeparator = "_"
)

var ErrInvalidPlatform = errors.New("invalid platform")

/* -------------------------------------------------------------------------- */
/*                             Struct: Executable                             */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing a Godot executable file.
type Executable struct {
	version  version.Version
	platform platform.Platform
}

// Compile-time verifications that 'Executable' implements 'Artifact'.
var _ artifact.Artifact = Executable{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Executable' struct with the specified values.
//
// NOTE: This method is rather pointless, but the 'Versioned' and 'Platformed'
// interfaces conflict with the desired field names.
func New(v version.Version, p platform.Platform) Executable {
	return Executable{version: v, platform: p}
}

/* --------------------------- Function: ToArchive -------------------------- */

// A convenience method which returns a Godot executable archive 'Artifact'.
func (ex Executable) ToArchive() Archive {
	return Archive{Inner: ex}
}

/* ----------------------------- Impl: Artifact ----------------------------- */

// Returns the name of the Godot executable, given the specified 'Version' and
// 'Platform'.
//
// NOTE: Godot names its executables in the format 'Godot_<VERSION>_<PLATFORM>',
// with Windows executables getting an extra '.exe' extension. Both the version
// and platform identifiers are version-specific, but the overall naming scheme
// has not changed (as of v4.2).
func (ex Executable) Name() string {
	var name strings.Builder

	name.WriteString(namePrefix)
	name.WriteString(nameSeparator)

	name.WriteString(ex.version.String())
	name.WriteString(nameSeparator)

	platformIdentifier, err := platform.Format(ex.platform, ex.version)
	if err != nil {
		return ""
	}

	name.WriteString(platformIdentifier)

	if ex.platform.OS == platform.Windows {
		name.WriteString(".exe")
	}

	return name.String()
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (ex Executable) Version() version.Version {
	return ex.version
}

/* ---------------------------- Impl: Platformed ---------------------------- */

func (ex Executable) Platform() platform.Platform {
	return ex.platform
}

/* ----------------------------- Impl: Stringer ----------------------------- */

func (ex Executable) String() string {
	return ex.Name()
}
