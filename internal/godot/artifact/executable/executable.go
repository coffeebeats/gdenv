package executable

import (
	"errors"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	namePrefix    = "Godot"
	nameSeparator = "_"

	pathGodotAppExecutable = "Contents/MacOS/Godot"
)

var ErrInvalidPlatform = errors.New("invalid platform")

/* -------------------------------------------------------------------------- */
/*                             Struct: Executable                             */
/* -------------------------------------------------------------------------- */

type Folder = artifact.Folder[Executable]
type Archive = archive.Zip[Folder]

// An 'Artifact' representing the Godot application itself.
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

/* ----------------------------- Function: Path ----------------------------- */

// The path relative to the 'Artifact' that executes Godot. For Linux and
// Windows this will be the executable 'Artifact' itself. On macOS this will be
// a path within the app folder.
func (ex Executable) Path() string {
	if ex.platform.OS == platform.MacOS {
		return filepath.Join(ex.Name(), pathGodotAppExecutable)
	}

	return ex.Name()
}

/* ----------------------------- Impl: Artifact ----------------------------- */

// Returns the top-level name of the application artifact that should be stored
// by 'gdenv'.
//
// NOTE: On Linux and Windows this is simply the executable itself, but on
// macOS an app folder structure is shipped. That folder should be saved for
// macOS platforms. As such, 'Godot.app' is returned for the name.
func (ex Executable) Name() string {
	return Name(ex.version, ex.platform)
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
