package artifacttest

/* -------------------------------------------------------------------------- */
/*                              Struct: Artifact                              */
/* -------------------------------------------------------------------------- */

type Artifact struct {
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Artifact) Name() string {
	return "artifact"
}
