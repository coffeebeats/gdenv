package download

import (
	"context"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

type localExArchive = artifact.Local[executable.Archive]
type localExChecksums = artifact.Local[checksum.Executable]

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Executable downloads the Godot 'executable.Archive' for a specific version
// and platform and returns an 'artifact.Local' encapsulating the result.
func Executable(
	ctx context.Context,
	m mirror.Mirror,
	ex executable.Executable,
	out string,
) (localExArchive, error) {
	if err := checkIsDirectory(out); err != nil {
		return localExArchive{}, err
	}

	remote, err := m.ExecutableArchive(ex.Version(), ex.Platform())
	if err != nil {
		return localExArchive{}, err
	}

	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localExArchive{}, err
	}

	return localExArchive{
		Artifact: remote.Artifact,
		Path:     filepath.Join(out, remote.Artifact.Name()),
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                        Function: ExecutableChecksums                       */
/* -------------------------------------------------------------------------- */

// ExecutableChecksums downloads the Godot 'checksum.Source' file for a specific
// version and returns an 'artifact.Local' encapsulating the result.
func ExecutableChecksums(
	ctx context.Context,
	m mirror.Mirror,
	v version.Version,
	out string,
) (localExChecksums, error) {
	if err := checkIsDirectory(out); err != nil {
		return localExChecksums{}, err
	}

	remote, err := m.ExecutableArchiveChecksums(v)
	if err != nil {
		return localExChecksums{}, err
	}

	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localExChecksums{}, err
	}

	return localExChecksums{
		Artifact: remote.Artifact,
		Path:     filepath.Join(out, remote.Artifact.Name()),
	}, nil
}
