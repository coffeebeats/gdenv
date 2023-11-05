package source

import (
	"crypto/sha256"
	"hash"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Checksums                             */
/* -------------------------------------------------------------------------- */

// An 'Artifact' representing a Godot source archive checksums file.
type Checksums struct {
	version version.Version
}

// Returns a new 'Checksums' struct after validating the Godot version.
func NewChecksums(v version.Version) (Checksums, error) {
	var c Checksums

	if v.CompareNormal(versionSourceChecksumsSupported()) < 0 {
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
	return archive.TarXZ[Source]{Inner: New(c.version)}.Name() + ".sha256"
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
	return sha256.New()
}

/* -------------------------------------------------------------------------- */
/*                  Function: versionSourceChecksumsSupported                 */
/* -------------------------------------------------------------------------- */

// Returns the first version at which 'godot-<version>.tar.xa.sha256' files for
// source archives began being published.
func versionSourceChecksumsSupported() version.Version {
	return version.MustParse("3.0")
}
