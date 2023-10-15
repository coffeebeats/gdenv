package download

import (
	"context"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"golang.org/x/sync/errgroup"
)

type localSourceArchive = artifact.Local[source.Archive]
type localSourceChecksums = artifact.Local[checksum.Source]

/* -------------------------------------------------------------------------- */
/*                   Function: SourceWithChecksumValidation                   */
/* -------------------------------------------------------------------------- */

func SourceWithChecksumValidation(
	ctx context.Context,
	m mirror.Mirror,
	v version.Version,
	out string,
) (artifact.Local[source.Archive], error) {
	chSource, chChecksums := make(chan artifact.Local[source.Archive]), make(chan artifact.Local[checksum.Source])
	defer close(chSource)
	defer close(chChecksums)

	eg, ctxDownload := errgroup.WithContext(ctx)

	eg.Go(func() error {
		result, err := Source(ctxDownload, m, v, out)
		if err != nil {
			return err
		}

		select {
		case chSource <- result:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	eg.Go(func() error {
		result, err := SourceChecksums(ctxDownload, m, v, out)
		if err != nil {
			return err
		}

		select {
		case chChecksums <- result:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	sourceArchive, checksums := <-chSource, <-chChecksums

	if err := eg.Wait(); err != nil {
		return localSourceArchive{}, err
	}

	if err := checksum.Compare[source.Archive](ctx, sourceArchive, checksums); err != nil {
		return localSourceArchive{}, err
	}

	return sourceArchive, nil
}

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

	out = filepath.Join(out, remote.Artifact.Name())
	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localSourceArchive{}, err
	}

	return localSourceArchive{
		Artifact: remote.Artifact,
		Path:     out,
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

	out = filepath.Join(out, remote.Artifact.Name())
	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localSourceChecksums{}, err
	}

	return localSourceChecksums{
		Artifact: remote.Artifact,
		Path:     out,
	}, nil
}
