package archive

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                               Struct: Archive                              */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip' compressed Godot application/binary.
//
// NOTE: 'Archive' structs are currently only used for executables, which
// are exclusively compressed into '.zip' files. If this assumpted changes,
// for example if source code '.tar.xz' archives need to be supported, then
// this struct will need to be differentiated into archive types.
type Archive struct {
	Executable executable.Executable
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Archive) Name() string {
	return a.Executable.Name() + extensionZip
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (a Archive) Version() version.Version {
	return a.Executable.Version()
}

/* ---------------------------- Impl: Platformed ---------------------------- */

func (a Archive) Platform() platform.Platform {
	return a.Executable.Platform()
}

/* ----------------------------- Impl: Extracter ---------------------------- */
