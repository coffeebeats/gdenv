package artifact

/* -------------------------------------------------------------------------- */
/*                                Struct: Local                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is available on the local file system.
type Local[A Artifact] struct {
	Artifact A
	Path     string
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (l Local[A]) Name() string {
	return l.Artifact.Name()
}
