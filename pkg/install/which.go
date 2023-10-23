package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
)

var ErrNotInstalled = errors.New("version not installed")

/* -------------------------------------------------------------------------- */
/*                               Function: Which                              */
/* -------------------------------------------------------------------------- */

// Which returns the path to the cached Godot executable specified by the
// locally or globally pinned version.
func Which(ctx context.Context, storePath string, p platform.Platform, atPath string) (string, error) {
	v, err := pin.VersionAt(ctx, storePath, atPath)
	if err != nil {
		return "", err
	}

	ex := executable.New(v, p)

	ok, err := store.Has(storePath, ex)
	if err != nil {
		return "", err
	}

	if !ok {
		// TODO: Determine whether this should be an error.
		return "", fmt.Errorf("%w: %s", ErrNotInstalled, v)
	}

	return store.Executable(storePath, ex)
}
