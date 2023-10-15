package checksum

import (
	"context"
	"errors"
	"fmt"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"golang.org/x/sync/errgroup"
)

var (
	ErrChecksumMismatch     = errors.New("checksum does not match")
	ErrChecksumsUnsupported = errors.New("version precedes checksums")
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Checksums                            */
/* -------------------------------------------------------------------------- */

// An interface for an 'Artifact' representing a checksums file.
type Checksums[T artifact.Artifact] interface {
	artifact.Artifact
	artifact.Versioned

	// NOTE: This dummy method is defined in order to (i) restrict outside
	// implementers and (ii) ensure the correct 'Artifact' types are used during
	// checksum extraction.
	supports(T)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Compare                             */
/* -------------------------------------------------------------------------- */

// Compare takes a local executable archive and a local checksums file for
// executable archives and validates that the executable archive's checksum
// matches the expected value.
func Compare[T archive.Archive, U Checksums[T]](
	ctx context.Context,
	localArtifact artifact.Local[T],
	localChecksums artifact.Local[U],
) error {
	eg, ctx := errgroup.WithContext(ctx)

	got, want := make(chan string, 1), make(chan string, 1)
	defer close(got)
	defer close(want)

	eg.Go(func() error {
		value, err := Compute[T](ctx, localArtifact)
		if err != nil {
			return err
		}

		select {
		case got <- value:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	eg.Go(func() error {
		value, err := Extract[T](ctx, localChecksums, localArtifact.Artifact)
		if err != nil {
			return err
		}

		select {
		case want <- value:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if g, w := <-got, <-want; g != w {
		return fmt.Errorf("%w: %s (got) != %s (want)", ErrChecksumMismatch, g, w)
	}

	return eg.Wait()
}

/* -------------------------------------------------------------------------- */
/*                              Struct: checksums                             */
/* -------------------------------------------------------------------------- */

// A shared implementation of a checksums file 'Artifact'; this should be
// wrapped by user-facing types.
type checksums struct {
	version version.Version
}
