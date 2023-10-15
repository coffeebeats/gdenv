package download

import (
	"context"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

type localSourceArchive = artifact.Local[source.Archive]
type localSourceChecksums = artifact.Local[checksum.Source]

/* -------------------------------------------------------------------------- */
/*                              Function: Source                              */
/* -------------------------------------------------------------------------- */

// Source downloads the Godot 'source.Archive' for a specific version and
// returns an 'artifact.Local' encapsulating the result.
func Source(
	ctx context.Context,
	m mirror.Mirror,
	v version.Version,
	out string,
) (localSourceArchive, error) {
	if err := checkIsDirectory(out); err != nil {
		return localSourceArchive{}, err
	}

	remote, err := m.SourceArchive(v)
	if err != nil {
		return localSourceArchive{}, err
	}

	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localSourceArchive{}, err
	}

	return localSourceArchive{
		Artifact: remote.Artifact,
		Path:     filepath.Join(out, remote.Artifact.Name()),
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                          Function: SourceChecksums                         */
/* -------------------------------------------------------------------------- */

// SourceChecksums downloads the Godot 'checksum.Source' file for a specific
// version and returns an 'artifact.Local' encapsulating the result.
func SourceChecksums(
	ctx context.Context,
	m mirror.Mirror,
	v version.Version,
	out string,
) (localSourceChecksums, error) {
	if err := checkIsDirectory(out); err != nil {
		return localSourceChecksums{}, err
	}

	remote, err := m.SourceArchiveChecksums(v)
	if err != nil {
		return localSourceChecksums{}, err
	}

	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localSourceChecksums{}, err
	}

	return localSourceChecksums{
		Artifact: remote.Artifact,
		Path:     filepath.Join(out, remote.Artifact.Name()),
	}, nil
}
