package checksums

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const FilenameChecksums = "SHA512-SUMS.txt"

var ErrChecksumsUnsupported = errors.New("version precedes checksums")

/* -------------------------------------------------------------------------- */
/*                              Struct: Checksums                             */
/* -------------------------------------------------------------------------- */

type Checksums struct {
	version version.Version
}

// Compile-time verifications that 'Checksums' implements 'Artifact'.
var _ artifact.Artifact = Checksums{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

func New(v version.Version) (Checksums, error) {
	var c Checksums

	if v.CompareNormal(versionChecksumsSupported()) < 0 {
		return c, fmt.Errorf("%w: %s", ErrChecksumsUnsupported, v)
	}

	c.version = v

	return c, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (c Checksums) Name() string {
	return FilenameChecksums
}

/* ----------------------------- Impl: Versioned ---------------------------- */

func (c Checksums) Version() version.Version {
	return c.version
}

/* -------------------------------------------------------------------------- */
/*                     Function: versionChecksumsSupported                    */
/* -------------------------------------------------------------------------- */

// Returns the first version at which 'SHA512-SUMS.txt' began being published
// alongside other published artifacts.
func versionChecksumsSupported() version.Version {
	return version.MustParse("3.2.2")
}
