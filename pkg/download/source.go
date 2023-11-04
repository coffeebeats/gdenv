package download

import (
	"context"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"golang.org/x/sync/errgroup"
)

/* -------------------------------------------------------------------------- */
/*                   Function: SourceWithChecksumValidation                   */
/* -------------------------------------------------------------------------- */

// SourceWithChecksumValidation downloads a source code archive and validates
// that its checksum matches the published value.
func SourceWithChecksumValidation(
	ctx context.Context,
	v version.Version,
	out string,
) (artifact.Local[source.Archive], error) {
	chSource, chChecksums := make(chan artifact.Local[source.Archive]), make(chan artifact.Local[checksum.Source])
	defer close(chSource)
	defer close(chChecksums)

	eg, ctxDownload := errgroup.WithContext(ctx)

	eg.Go(func() error {
		srcArchive := source.Archive{Artifact: source.New(v)}

		result, err := Download(ctxDownload, srcArchive, out)
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
		checksums, err := checksum.NewSource(v)
		if err != nil {
			return err
		}

		result, err := Download(ctxDownload, checksums, out)
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
		return artifact.Local[source.Archive]{}, err
	}

	if err := checksum.Compare[source.Archive](ctx, sourceArchive, checksums); err != nil {
		return artifact.Local[source.Archive]{}, err
	}

	return sourceArchive, nil
}
