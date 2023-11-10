package download

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
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
	chArchive := make(chan artifact.Local[source.Archive])
	defer close(chArchive)

	chChecksums := make(chan artifact.Local[source.Checksums])
	defer close(chChecksums)

	eg, ctxDownload := errgroup.WithContext(ctx)

	eg.Go(func() error {
		srcArchive := source.Archive{Inner: source.New(v)}

		result, err := Download(ctxDownload, srcArchive, out)
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
		checksums, err := source.NewChecksums(v)
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

	srcArchive, checksums := <-chArchive, <-chChecksums

	if err := eg.Wait(); err != nil {
		return artifact.Local[source.Archive]{}, err
	}

	if err := checksum.Compare[source.Archive](ctx, srcArchive, checksums); err != nil {
		return artifact.Local[source.Archive]{}, err
	}

	return srcArchive, nil
}
