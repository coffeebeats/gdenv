package download

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/progress"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/mirror"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"golang.org/x/sync/errgroup"
)

type (
	progressKeyExecutable         struct{}
	progressKeyExecutableChecksum struct{}

	localExArchive   = artifact.Local[executable.Archive]
	localExChecksums = artifact.Local[checksum.Executable]
)

/* -------------------------------------------------------------------------- */
/*                      Function: WithExecutableProgress                      */
/* -------------------------------------------------------------------------- */

// WithSourceProgress creates a sub-context with an associated progress
// reporter. The result can be passed to download functions in this package to
// get updates on download progress.
func WithExecutableProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKeyExecutable{}, p)
}

/* -------------------------------------------------------------------------- */
/*                  Function: WithExecutableChecksumProgress                  */
/* -------------------------------------------------------------------------- */

// WithExecutableChecksumProgress creates a sub-context with an associated
// progress reporter. The result can be passed to download functions in this
// package to get updates on download progress.
func WithExecutableChecksumProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKeyExecutableChecksum{}, p)
}

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

	eg, downloadCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		result, err := Executable(downloadCtx, m, ex, out)
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
		result, err := ExecutableChecksums(downloadCtx, m, ex.Version(), out)
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

	executableMirror, ok := m.(mirror.Executable)
	if !ok || executableMirror == nil {
		return localExArchive{}, fmt.Errorf("%w: executables", mirror.ErrNotSupported)
	}

	remote, err := executableMirror.ExecutableArchive(ex.Version(), ex.Platform())
	if err != nil {
		return localExArchive{}, err
	}

	c := client.NewWithRedirectDomains(m.Domains()...)

	out = filepath.Join(out, remote.Artifact.Name())
	if err := downloadArtifact(ctx, c, remote, out, progressKeyExecutable{}); err != nil {
		return localExArchive{}, err
	}

	log.Debugf("downloaded executable: %s", out)

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

	executableMirror, ok := m.(mirror.Executable)
	if !ok || executableMirror == nil {
		return localExChecksums{}, fmt.Errorf("%w: executables", mirror.ErrNotSupported)
	}

	remote, err := executableMirror.ExecutableArchiveChecksums(v)
	if err != nil {
		return localExChecksums{}, err
	}

	c := client.NewWithRedirectDomains(m.Domains()...)

	out = filepath.Join(out, remote.Artifact.Name())
	if err := downloadArtifact(ctx, c, remote, out, progressKeyExecutableChecksum{}); err != nil {
		return localExChecksums{}, err
	}

	log.Debugf("downloaded checksums file: %s", out)

	return localExChecksums{
		Artifact: remote.Artifact,
		Path:     out,
	}, nil
}
