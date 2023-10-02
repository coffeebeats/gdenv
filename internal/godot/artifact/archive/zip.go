package archive

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip'-compressed archive.
type Zip[T Archivable] struct {
	Artifact T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Zip[T]) Name() string {
	name := a.Artifact.Name()
	if name != "" {
		name += extensionZip
	}

	return name
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
func (a Zip[T]) extract(path, out string) error { //nolint:revive
	return nil // TODO: Implement the archive extraction.
}
