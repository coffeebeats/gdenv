package install

import (
	"context"
	"errors"
	"io/fs"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
)

var ErrGodotNotFound = errors.New("godot not found")

/* -------------------------------------------------------------------------- */
/*                               Function: Which                              */
/* -------------------------------------------------------------------------- */

// Which returns the path to the cached Godot executable specified by the
// locally or globally pinned version.
func Which(ctx context.Context, storePath string, p platform.Platform, atPath string) (string, error) {
	path, err := pin.Resolve(ctx, atPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
	}

	// No pin file was found yet, so check globally.
	if path == "" {
		path = storePath
	}

	v, err := pin.Read(path)
	if err != nil {
		return "", ErrGodotNotFound
	}

	// Define the target 'Executable'.
	ex := executable.New(v, p)

	if !store.Has(storePath, ex) {
		return "", ErrGodotNotFound
	}

	return store.ToolPath(storePath, ex)
}
