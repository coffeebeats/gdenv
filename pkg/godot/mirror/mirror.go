package mirror

import (
	"context"
	"errors"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

var (
	ErrInvalidURL          = errors.New("invalid URL")
	ErrMissingMirrors      = errors.New("no mirrors provided")
	ErrNotFound            = errors.New("no mirror found")
	ErrUnsupportedArtifact = errors.New("unsupported artifact")
)

// clientKey is a context key used internally to replace the REST client used.
type clientKey struct{}

/* -------------------------------------------------------------------------- */
/*                              Interface: Mirror                             */
/* -------------------------------------------------------------------------- */

// Mirror specifies a host of Godot release artifacts.
type Mirror[T artifact.Artifact] interface {
	Hoster
	Remoter[T]

	Name() string
}

/* -------------------------------------------------------------------------- */
/*                              Interface: Hoster                             */
/* -------------------------------------------------------------------------- */

// Hoster is a mirror which describes the host URLs at which it hosts content.
// This can be used to restrict redirects when downloading artifacts, improving
// security.
type Hoster interface {
	// Hosts returns a slice of URL hosts at which the mirror hosts artifacts.
	Hosts() []string
}

/* -------------------------------------------------------------------------- */
/*                             Interface: Remoter                             */
/* -------------------------------------------------------------------------- */

// Remoter is a type that can resolve the URL at which a specified artifact is
// hosted. Provided artifacts must be versioned.
type Remoter[T artifact.Artifact] interface {
	Remote(a T) (artifact.Remote[T], error)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Select                              */
/* -------------------------------------------------------------------------- */

// Select chooses the best 'Mirror' of those provided for downloading the
// specified Godot release artifact.
func Select[T artifact.Artifact](
	ctx context.Context,
	mirrors []Mirror[T],
	a T,
) (Mirror[T], error) {
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

	selected := make(chan Mirror[T])

	for _, m := range mirrors {
		m := m // Prevent capture of loop variable.

		eg.Go(func() error {
			ok, err := checkIfExists(ctx, m, a)
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
func checkIfExists[T artifact.Artifact](
	ctx context.Context,
	m Mirror[T],
	a T,
) (bool, error) {
	remote, err := m.Remote(a)
	if err != nil {
		return false, err
	}

	// NOTE: It would be cleaner to expose this as an actual dependency, as an
	// HTTP client *is* required. However, the internal 'client.Client'
	// implementation is opinionated and not ready to be exposed yet as a public
	// type. For now, this simply allows tests to inject a client.
	c, ok := ctx.Value(clientKey{}).(*client.Client)
	if !ok || c == nil {
		c = client.NewWithRedirectDomains(m.Hosts()...)
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
func chooseBest[T artifact.Artifact](
	available <-chan Mirror[T],
	ranking []Mirror[T],
) (Mirror[T], error) {
	out, index := Mirror[T](nil), len(ranking)

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
