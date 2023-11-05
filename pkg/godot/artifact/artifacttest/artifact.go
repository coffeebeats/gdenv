package artifacttest

import "github.com/coffeebeats/gdenv/pkg/godot/version"

/* -------------------------------------------------------------------------- */
/*                            Struct: MockArtifact                            */
/* -------------------------------------------------------------------------- */

type MockArtifact struct {
	name    string
	version version.Version
}

/* ------------------------ Function: NewWithVersion ------------------------ */

// NewWithVersion creates a new mock artifact with the specified version.
func NewWithVersion(v version.Version) MockArtifact {
	return MockArtifact{name: "", version: v}
}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'MockArtifact' as a Godot release artifact.
func (a MockArtifact) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

func (a MockArtifact) Name() string {
	return a.name
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (a MockArtifact) Version() version.Version {
	return a.version
}

/* ------------------------ Impl: archive.Archivable ------------------------ */

func (a MockArtifact) Archivable() {}
