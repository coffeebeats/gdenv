package executable

import (
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	namePrefix    = "Godot"
	nameSeparator = "_"

	nameGodotMacOSApp = "Godot.app"
)

type Archive = archive.Zip[Executable]

/* -------------------------------------------------------------------------- */
/*                             Struct: Executable                             */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing the Godot application itself.
type Executable struct {
	version  version.Version
	platform platform.Platform
}

// Compile-time verifications that 'Executable' implements 'Artifact'.
var _ artifact.Artifact = (*Executable)(nil)

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Executable' struct with the specified values.
//
// NOTE: This method is rather pointless, but the 'Versioned' and 'Platformed'
// interfaces conflict with the desired field names.
func New(v version.Version, p platform.Platform) Executable {
	return Executable{version: v, platform: p}
}

/* ----------------------------- Function: Path ----------------------------- */

// The path relative to the 'Artifact' that executes Godot. For Linux and
// Windows this will be the executable 'Artifact' itself. On macOS this will be
// a path within the app folder.
func (ex Executable) Path() string {
	if ex.platform.OS == platform.MacOS {
		return filepath.Join(nameGodotMacOSApp, pathGodotAppExecutable())
	}

	return ex.Name()
}

/* ------------------------ Impl: archive.Archivable ------------------------ */

// Allows 'Executable' to be used by 'Archive' implementation.
func (ex Executable) Archivable() {}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'Executable' as a Godot release artifact.
func (ex Executable) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

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

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (ex Executable) Version() version.Version {
	return ex.version
}

/* ------------------------ Impl: artifact.Platformed ----------------------- */

func (ex Executable) Platform() platform.Platform {
	return ex.platform
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (ex Executable) String() string {
	return ex.Name()
}

/* -------------------------------------------------------------------------- */
/*                      Function: pathGodotAppExecutable                      */
/* -------------------------------------------------------------------------- */

// Returns an OS-specific path segment from the macOS Godot application to the
// executable file.
func pathGodotAppExecutable() string {
	return filepath.Join("Contents", "MacOS", "Godot")
}
