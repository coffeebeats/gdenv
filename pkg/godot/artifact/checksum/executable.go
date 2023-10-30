package checksum

import (
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const filenameChecksums = "SHA512-SUMS.txt"

/* -------------------------------------------------------------------------- */
/*                             Struct: Executable                             */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing a Godot source archive checksums file.
type Executable checksums

// Compile-time verifications that 'Executable' implements 'Artifact'.
var _ artifact.Artifact = Executable{} //nolint:exhaustruct

// Compile-time verifications that 'Executable' implements 'Checksums'.
var _ Checksums[executable.Archive] = Executable{} //nolint:exhaustruct

/* ------------------------- Function: NewExecutable ------------------------ */

// Returns a new 'Executable' struct after validating the Godot version.
func NewExecutable(v version.Version) (Executable, error) {
	var ex Executable

	if v.CompareNormal(versionExecutableChecksumsSupported()) < 0 {
		return ex, ErrChecksumsUnsupported
	}

	ex.version = v

	return ex, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (ex Executable) Name() string {
	return filenameChecksums
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (ex Executable) Version() version.Version {
	return ex.version
}

/* ----------------------------- Impl: Checksums ---------------------------- */

func (ex Executable) supports(executable.Archive) {} //nolint:unused

/* -------------------------------------------------------------------------- */
/*                Function: versionExecutableChecksumsSupported               */
/* -------------------------------------------------------------------------- */

// Returns the first version at which 'SHA512-SUMS.txt' files for executable
// archives began being published.
func versionExecutableChecksumsSupported() version.Version {
	return version.MustParse("3.2.2")
}
