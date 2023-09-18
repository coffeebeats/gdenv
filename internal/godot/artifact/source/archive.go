package source

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
)

type Archive = archive.TarXZ[Source]

// Compile-time verifications that 'Archive' implements required interfaces.
var _ artifact.Artifact = Archive{} //nolint:exhaustruct
