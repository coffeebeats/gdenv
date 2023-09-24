package artifact

/* -------------------------------------------------------------------------- */
/*                               Struct: Folder                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is a directory containing another
// 'Artifact'.
type Folder[T Artifact] struct {
	Inner      T
	FolderName string
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (f Folder[T]) Name() string {
	return f.FolderName
}

/* ------------------------------ Impl: Wrapper ----------------------------- */

func (f Folder[T]) Contents() T { //nolint:ireturn
	return f.Inner
}
