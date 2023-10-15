package download

import (
	"context"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"golang.org/x/sync/errgroup"
)

type localExArchive = artifact.Local[executable.Archive]
type localExChecksums = artifact.Local[checksum.Executable]

/* -------------------------------------------------------------------------- */
/*                 Function: ExecutableWithChecksumValidation                 */
/* -------------------------------------------------------------------------- */

func ExecutableWithChecksumValidation(
	ctx context.Context,
	m mirror.Mirror,
	ex executable.Executable,
	out string,
) (artifact.Local[executable.Archive], error) {
	chArchive, chChecksums := make(chan artifact.Local[executable.Archive]), make(chan artifact.Local[checksum.Executable])
	defer close(chArchive)
	defer close(chChecksums)

	eg, ctxEG := errgroup.WithContext(ctx)

	eg.Go(func() error {
		result, err := Executable(ctxEG, m, ex, out)
		if err != nil {
			return err
		}

		select {
		case chArchive <- result:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	eg.Go(func() error {
		result, err := ExecutableChecksums(ctxEG, m, ex.Version(), out)
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

	exArchive, exArchiveChecksums := <-chArchive, <-chChecksums

	if err := eg.Wait(); err != nil {
		return artifact.Local[executable.Archive]{}, err
	}

	if err := checksum.Compare[executable.Archive](ctx, exArchive, exArchiveChecksums); err != nil {
		return artifact.Local[executable.Archive]{}, err
	}

	return exArchive, nil
}

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

	out = filepath.Join(out, remote.Artifact.Name())
	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localExArchive{}, err
	}

	return localExArchive{
		Artifact: remote.Artifact,
		Path:     out,
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

	out = filepath.Join(out, remote.Artifact.Name())
	if err := m.Client().DownloadTo(ctx, remote.URL, out); err != nil {
		return localExChecksums{}, err
	}

	return localExChecksums{
		Artifact: remote.Artifact,
		Path:     out,
	}, nil
}
