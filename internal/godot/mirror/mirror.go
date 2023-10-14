package mirror

import (
	"context"
	"errors"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidSpecification = errors.New("invalid specification")
	ErrInvalidURL           = errors.New("invalid URL")
	ErrNotFound             = errors.New("no mirror found")
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// An interface specifying methods for retrieving information about assets
// available for download via a mirror host.
type Mirror interface {
	Client() client.Client
	ExecutableArchive(version.Version, platform.Platform) (artifact.Remote[executable.Archive], error)
	ExecutableArchiveChecksums(version.Version) (artifact.Remote[checksum.Executable], error)

	SourceArchive(version.Version) (artifact.Remote[source.Archive], error)
	SourceArchiveChecksums(version.Version) (artifact.Remote[checksum.Source], error)

	// Issues a request to see if the mirror host has the specific version.
	CheckIfExists(context.Context, version.Version) bool

	// Checks whether the version is broadly supported by the mirror. No network
	// request is issued, but this does not guarantee the host has the version.
	// To check whether the host has the version definitively via the network,
	// use the 'Has' method.
	Supports(version.Version) bool
}

/* -------------------------------------------------------------------------- */
/*                              Function: Choose                              */
/* -------------------------------------------------------------------------- */

// Choose selects the best 'Mirror' for downloading assets for the specified
// version of Godot.
func Choose(ctx context.Context, v version.Version) (Mirror, error) { //nolint:cyclop,funlen,ireturn
	eg, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()

		// NOTE: The 'errgroup.Group' needs to be waited to prevent a leak.
		eg.Wait() //nolint:errcheck
	}()

	selected := make(chan Mirror)

	// Check if 'GitHub' supports the specified version.
	eg.Go(func() error {
		// NOTE: Use a zero value to avoid initializing a client before necessary.
		if !(GitHub{}).Supports(v) { //nolint:exhaustruct
			return nil
		}

		m := NewGitHub()
		if !m.CheckIfExists(ctx, v) {
			return nil
		}

		select {
		case selected <- m:
		case <-ctx.Done():
			break
		}

		return nil
	})

	// Check if 'TuxFamily' supports the specified version.
	eg.Go(func() error {
		// NOTE: Use a zero value to avoid initializing a client before necessary.
		if !(TuxFamily{}).Supports(v) { //nolint:exhaustruct
			return nil
		}

		m := NewTuxFamily()
		if !m.CheckIfExists(ctx, v) {
			return nil
		}

		select {
		case selected <- m:
		case <-ctx.Done():
			break
		}

		return nil
	})

	go func() {
		eg.Wait() //nolint:errcheck
		close(selected)
	}()

	var out Mirror

	for m := range selected {
		// Take GitHub immediately if it's valid.
		if m, ok := m.(GitHub); ok {
			return m, nil
		}

		if out == nil {
			out = m
			continue
		}
	}

	return out, eg.Wait()
}
