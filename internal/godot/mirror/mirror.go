package mirror

import (
	"context"
	"errors"
	"fmt"

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
	ErrNotSupported         = errors.New("mirror not supported")
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// Specifies a host of Godot release artifacts. The associated methods are
// related to the host itself and not about individual artifacts.
type Mirror interface {
	// Domains returns a slice of domains at which the mirror hosts artifacts.
	Domains() []string

	// Checks whether the version is broadly supported by the mirror. No network
	// request is issued, but this does not guarantee the host has the version.
	// To check whether the host has the version definitively via the network,
	// use the 'Has' method.
	Supports(v version.Version) bool
}

/* -------------------------------------------------------------------------- */
/*                            Interface: Executable                           */
/* -------------------------------------------------------------------------- */

// Executable is a mirror which hosts Godot executable artifacts. This does not
// imply that *all* executable versions are hosted, so users should be prepared
// to handle the case where resolving the artifact URL fails.
type Executable interface {
	Mirror

	ExecutableArchive(v version.Version, p platform.Platform) (artifact.Remote[executable.Archive], error)
	ExecutableArchiveChecksums(v version.Version) (artifact.Remote[checksum.Executable], error)
}

/* -------------------------------------------------------------------------- */
/*                              Interface: Source                             */
/* -------------------------------------------------------------------------- */

// Source is a mirror which hosts Godot repository source code versions. This
// does not imply that *all* executable versions are hosted, so users should be
// prepared to handle the case where resolving the artifact URL fails.
type Source interface {
	Mirror

	SourceArchive(v version.Version) (artifact.Remote[source.Archive], error)
	SourceArchiveChecksums(v version.Version) (artifact.Remote[checksum.Source], error)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Choose                              */
/* -------------------------------------------------------------------------- */

// Choose selects the best 'Mirror' for downloading assets for the specified
// version and platform of Godot.
func Choose( //nolint:cyclop,funlen,ireturn
	ctx context.Context,
	v version.Version,
	p platform.Platform,
) (Mirror, error) {
	eg, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()

		// NOTE: Wait on 'errgroup.Group'to prevent goroutine leaks.
		eg.Wait() //nolint:errcheck
	}()

	selected := make(chan Mirror)

	// Check if 'GitHub' supports the specified version.
	eg.Go(func() error {
		// NOTE: Use a zero value to avoid initializing a client before necessary.
		if !(GitHub{}).Supports(v) {
			return nil
		}

		m := GitHub{}
		ok, err := checkIfExists(ctx, m, v, p)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		select {
		case selected <- m:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	// Check if 'TuxFamily' supports the specified version.
	eg.Go(func() error {
		// NOTE: Use a zero value to avoid initializing a client before necessary.
		if !(TuxFamily{}).Supports(v) {
			return nil
		}

		m := TuxFamily{}
		ok, err := checkIfExists(ctx, m, v, p)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		select {
		case selected <- m:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	go func() {
		eg.Wait() //nolint:errcheck
		close(selected)
	}()

	var out Mirror

	for m := range selected {
		// Take GitHub immediately if it's valid because it has great speeds.
		if m, ok := m.(GitHub); ok {
			return m, nil
		}

		out = m
	}

	if out == nil {
		return nil, fmt.Errorf("%w: version '%s'", ErrNotFound, v)
	}

	return out, eg.Wait()
}

/* ------------------------- Function: checkIfExists ------------------------ */

// Issues a request to the mirror host to determine if the artifact exists.
func checkIfExists(
	ctx context.Context,
	m Executable,
	v version.Version,
	p platform.Platform,
) (bool, error) {
	if !m.Supports(v) {
		return false, nil
	}

	remote, err := m.ExecutableArchive(v, p)
	if err != nil {
		return false, err
	}

	c := client.NewWithRedirectDomains(m.Domains()...)

	exists, err := c.Exists(ctx, remote.URL.String())
	if err != nil {
		return false, err
	}

	return exists, nil
}
