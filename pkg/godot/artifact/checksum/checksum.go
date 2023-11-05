package checksum

import (
	"context"
	"errors"
	"fmt"
	"hash"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
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

	// Supports is a method to register which artifact the checksums are for.
	Supports(_ T)

	// Hash returns a 'hash.Hash' instance used to compute the checksum of the
	// supported artifact type.
	Hash() hash.Hash
}

/* -------------------------------------------------------------------------- */
/*                              Function: Compare                             */
/* -------------------------------------------------------------------------- */

// Compare takes a local executable archive and a local checksums file for
// executable archives and validates that the executable archive's checksum
// matches the expected value.
func Compare[T artifact.Artifact, U Checksums[T]](
	ctx context.Context,
	localArtifact artifact.Local[T],
	localChecksums artifact.Local[U],
) error {
	log.Info("verifying checksum of downloaded file")

	eg, ctx := errgroup.WithContext(ctx)

	got, want := make(chan string, 1), make(chan string, 1)
	defer close(got)
	defer close(want)

	eg.Go(func() error {
		value, err := Compute[T](ctx, localChecksums.Artifact.Hash(), localArtifact)
		if err != nil {
			return err
		}

		log.Debugf("actual checksum: %s", value)

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

		log.Debugf("expected checksum: %s", value)

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

	log.Debug("checksum matched expected value")

	return eg.Wait()
}
