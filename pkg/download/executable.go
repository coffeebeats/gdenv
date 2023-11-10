package download

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
)

/* -------------------------------------------------------------------------- */
/*                 Function: ExecutableWithChecksumValidation                 */
/* -------------------------------------------------------------------------- */

// ExecutableWithChecksumValidation downloads an executable archive and
// validates that its checksum matches the published value.
func ExecutableWithChecksumValidation(
	ctx context.Context,
	ex executable.Executable,
	out string,
) (artifact.Local[executable.Archive], error) {
	chArchive := make(chan artifact.Local[executable.Archive])
	defer close(chArchive)

	chChecksums := make(chan artifact.Local[executable.Checksums])
	defer close(chChecksums)

	eg, ctxDownload := errgroup.WithContext(ctx)

	eg.Go(func() error {
		exArchive := executable.Archive{Inner: ex}

		result, err := Download(ctxDownload, exArchive, out)
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
		checksums, err := executable.NewChecksums(ex.Version())
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

	exArchive, exArchiveChecksums := <-chArchive, <-chChecksums

	if err := eg.Wait(); err != nil {
		return artifact.Local[executable.Archive]{}, err
	}

	if err := checksum.Compare[executable.Archive](ctx, exArchive, exArchiveChecksums); err != nil {
		return artifact.Local[executable.Archive]{}, err
	}

	return exArchive, nil
}
