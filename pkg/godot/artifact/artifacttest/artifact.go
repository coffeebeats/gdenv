package artifacttest

import "github.com/coffeebeats/gdenv/pkg/godot/version"

/* -------------------------------------------------------------------------- */
/*                            Struct: MockArtifact                            */
/* -------------------------------------------------------------------------- */

type MockArtifact struct {
	name    string
	version version.Version
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a MockArtifact) Name() string {
	return a.name
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (a MockArtifact) Version() version.Version {
	return a.version
}

/* ---------------------------- Impl: Archivable ---------------------------- */

func (a MockArtifact) Archivable() {}
