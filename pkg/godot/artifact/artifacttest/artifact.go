package artifacttest

/* -------------------------------------------------------------------------- */
/*                            Struct: MockArtifact                            */
/* -------------------------------------------------------------------------- */

type MockArtifact struct {
	name string
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a MockArtifact) Name() string {
	return a.name
}

/* ---------------------------- Impl: Archivable ---------------------------- */

func (a MockArtifact) Archivable() {}
