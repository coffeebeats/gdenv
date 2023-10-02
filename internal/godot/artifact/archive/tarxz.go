package archive

const extensionTarXZ = ".tar.xz"

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

// A struct representing an 'XZ'-compressed archive.
type TarXZ[T Archivable] struct {
	Artifact T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a TarXZ[T]) Name() string {
	name := a.Artifact.Name()
	if name != "" {
		name += extensionTarXZ
	}

	return name
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
func (a TarXZ[T]) extract(path, out string) error { //nolint:revive
	return nil // TODO: Implement the archive extraction.
}
