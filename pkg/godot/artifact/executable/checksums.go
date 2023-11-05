package executable

import (
	"crypto/sha512"
	"errors"
	"hash"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const filenameChecksums = "SHA512-SUMS.txt"

var ErrChecksumsUnsupported = errors.New("version precedes checksums")

/* -------------------------------------------------------------------------- */
/*                              Struct: Checksums                             */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing a Godot executable archive checksums file.
type Checksums struct {
	version version.Version
}

// Returns a new 'Checksums' struct after validating the Godot version.
func NewChecksums(v version.Version) (Checksums, error) {
	var c Checksums

	if v.CompareNormal(versionExecutableChecksumsSupported()) < 0 {
		return c, checksum.ErrChecksumsUnsupported
	}

	c.version = v

	return c, nil
}

/* ------------------------- Impl: artifact.Artifact ------------------------ */

// Artifact "registers" 'Checksums' as a Godot release artifact.
func (c Checksums) Artifact() {}

/* -------------------------- Impl: artifact.Named -------------------------- */

func (c Checksums) Name() string {
	return filenameChecksums
}

/* ------------------------ Impl: artifact.Versioned ------------------------ */

func (c Checksums) Version() version.Version {
	return c.version
}

/* ------------------------ Impl: checksum.Checksums ------------------------ */

// Supports "registers" 'Checksums' as containing checksums for the specified
// artifact type.
func (c Checksums) Supports(_ Archive) {}

// Hash returns a new 'hash.Hash' for computing the file hash of an executable.
func (c Checksums) Hash() hash.Hash {
	return sha512.New()
}

/* -------------------------------------------------------------------------- */
/*                Function: versionExecutableChecksumsSupported               */
/* -------------------------------------------------------------------------- */

// Returns the first version at which 'SHA512-SUMS.txt' files for executable
// archives began being published.
func versionExecutableChecksumsSupported() version.Version {
	return version.MustParse("3.2.2")
}
