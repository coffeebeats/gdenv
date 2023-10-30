package mirror

import (
	"context"
	"errors"
	"slices"

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
	ErrMissingMirrors       = errors.New("no mirrors provided")
	ErrNotFound             = errors.New("no mirror found")
	ErrNotSupported         = errors.New("mirror not supported")
)

// clientKey is a context key used internally to replace the REST client used.
type clientKey struct{}

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
/*                              Function: Select                              */
/* -------------------------------------------------------------------------- */

// Select chooses the best 'Mirror' of those provided for downloading assets
// corresponding to the specified version and platform of Godot.
func Select( //nolint:ireturn
	ctx context.Context,
	v version.Version,
	p platform.Platform,
	mirrors []Mirror,
) (Mirror, error) {
	if len(mirrors) == 0 {
		return nil, ErrMissingMirrors
	}

	eg, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()

		// NOTE: Wait on 'errgroup.Group'to prevent goroutine leaks.
		eg.Wait() //nolint:errcheck
	}()

	selected := make(chan Mirror)

	for _, m := range mirrors {
		executableMirror, ok := m.(Executable)
		if !ok || executableMirror == nil {
			continue
		}

		eg.Go(func() error {
			ok, err := checkIfExists(ctx, executableMirror, v, p)
			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			select {
			case selected <- executableMirror:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})
	}

	go func() {
		eg.Wait() //nolint:errcheck
		close(selected)
	}()

	out, err := chooseBest(selected, mirrors)
	if err != nil {
		return nil, err
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

	// NOTE: It would be cleaner to expose this as an actual dependency, as an
	// HTTP client *is* required. However, the internal 'client.Client'
	// implementation is opinionated and not ready to be exposed yet as a public
	// type. For now, this simply allows tests to inject a client.
	c, ok := ctx.Value(clientKey{}).(*client.Client)
	if !ok || c == nil {
		c = client.NewWithRedirectDomains(m.Domains()...)
	}

	exists, err := c.Exists(ctx, remote.URL.String())
	if err != nil {
		return false, err
	}

	return exists, nil
}

/* -------------------------- Function: chooseBest -------------------------- */

// chooseBest selects the best mirror from those available. The lowest indexed
// 'Mirror' in 'ranking' will be returned. If none are available an error is
// returned.
func chooseBest(available <-chan Mirror, ranking []Mirror) (Mirror, error) { //nolint:ireturn
	out, index := Mirror(nil), len(ranking)

	for m := range available {
		// Rank mirrors according to order in 'mirrors'.
		i := slices.Index(ranking, m)
		if i <= index {
			out = m
			index = i
		}
	}

	if out == nil {
		return nil, ErrNotFound
	}

	return out, nil
}
