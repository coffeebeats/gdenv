package pin

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

var (
	ErrMissingPin     = errors.New("missing version pin")
	ErrUnexpectedFile = errors.New("unexpected file")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Read                               */
/* -------------------------------------------------------------------------- */

// Parses a 'Version' from the specified pin file.
func Read(path string) (version.Version, error) {
	if path == "" {
		return version.Version{}, ErrMissingPath
	}

	path, err := clean(path)
	if err != nil {
		return version.Version{}, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return version.Version{}, err
		}

		return version.Version{}, fmt.Errorf("%w: '%s'", ErrMissingPin, path)
	}

	v, err := version.Parse(string(bytes))
	if err != nil {
		return version.Version{}, err
	}

	return v, nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: Remove                              */
/* -------------------------------------------------------------------------- */

// Deletes the specified pin file if it exists.
func Remove(path string) error {
	path, err := clean(path)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                             Function: VersionAt                            */
/* -------------------------------------------------------------------------- */

// Resolves a version for the specified directory. This function starts by
// looking for a pin file in the specified directory or any ancestor
// directories. If none are found then globally-pinned version is checked.
func VersionAt(ctx context.Context, storePath, path string) (version.Version, error) {
	path, err := clean(path)
	if err != nil {
		return version.Version{}, err
	}

	path = filepath.Dir(path)
	root := filepath.VolumeName(path) + string(os.PathSeparator)

	// Check if the specified path (or any ancestors) has a pin
	for path != root {
		if ctx.Err() != nil {
			return version.Version{}, ctx.Err()
		}

		pinPath := filepath.Join(path, pinFilename)

		info, err := os.Stat(pinPath)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return version.Version{}, err
			}

			path = filepath.Dir(path)

			continue
		}

		// Validate that the file is a regular file; this catches cases where
		// there's a directory named after 'pinFilename'.
		if !info.Mode().IsRegular() {
			return version.Version{}, fmt.Errorf("%w: '%s'", fs.ErrInvalid, pinPath)
		}

		break
	}

	// Try reading a global pin file if the specified directory and all
	// ancestors were missing pin files.
	if path == root {
		path = storePath
	}

	return Read(path)
}

/* -------------------------------------------------------------------------- */
/*                               Function: Write                              */
/* -------------------------------------------------------------------------- */

// Writes a 'Version' to the specified pin file path.
//
// NOTE: This function will fail if any directories along the path do not exist.
func Write(v version.Version, path string) error {
	path, err := clean(path)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(v.String()), osutil.ModeUserRW)
}
