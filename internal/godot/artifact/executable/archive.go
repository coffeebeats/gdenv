package executable

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
)

// A convenience wrapper around 'Zip[Executable]'.
type Archive = archive.Zip[artifact.Folder[Executable]]

// Compile-time verifications that 'Archive' implements required interfaces.
var _ artifact.Artifact = Archive{} //nolint:exhaustruct
