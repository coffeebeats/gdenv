package download

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/mirror"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/progress"
	"golang.org/x/sync/errgroup"
)

type (
	progressKeySource         struct{}
	progressKeySourceChecksum struct{}

	localSourceArchive   = artifact.Local[source.Archive]
	localSourceChecksums = artifact.Local[checksum.Source]
)

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithSourceProgress creates a sub-context with an associated progress
// reporter. The result can be passed to download functions in this package to
// get updates on download progress.
func WithSourceProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKeySource{}, p)
}

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithSourceChecksumProgress creates a sub-context with an associated progress
// reporter. The result can be passed to download functions in this package to
// get updates on download progress.
func WithSourceChecksumProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKeySourceChecksum{}, p)
}

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

	sourceMirror, ok := m.(mirror.Source)
	if !ok || sourceMirror == nil {
		return localSourceArchive{}, fmt.Errorf("%w: source code", mirror.ErrNotSupported)
	}

	remote, err := sourceMirror.SourceArchive(v)
	if err != nil {
		return localSourceArchive{}, err
	}

	c := client.NewWithRedirectDomains(m.Domains()...)

	out = filepath.Join(out, remote.Artifact.Name())
	if err := downloadArtifact(ctx, c, remote, out, progressKeySource{}); err != nil {
		return localSourceArchive{}, err
	}

	log.Debugf("downloaded source: %s", out)

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

	sourceMirror, ok := m.(mirror.Source)
	if !ok || sourceMirror == nil {
		return localSourceChecksums{}, fmt.Errorf("%w: source code", mirror.ErrNotSupported)
	}

	remote, err := sourceMirror.SourceArchiveChecksums(v)
	if err != nil {
		return localSourceChecksums{}, err
	}

	c := client.NewWithRedirectDomains(m.Domains()...)

	out = filepath.Join(out, remote.Artifact.Name())
	if err := downloadArtifact(ctx, c, remote, out, progressKeySourceChecksum{}); err != nil {
		return localSourceChecksums{}, err
	}

	log.Debugf("downloaded checksums file: %s", out)

	return localSourceChecksums{
		Artifact: remote.Artifact,
		Path:     out,
	}, nil
}
