package pin

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/pkg/godot"
)

var (
	ErrParseVersion   = errors.New("failed to parse version")
	ErrUnexpectedFile = errors.New("unexpected file")
)

/* ----------------------------- Function: Read ----------------------------- */

// Parses a 'Version' from the specified pin file.
func Read(path string) (godot.Version, error) {
	path, err := Clean(path)
	if err != nil {
		return godot.Version{}, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return godot.Version{}, err
	}

	version, err := godot.ParseVersion(string(bytes))
	if err != nil {
		return godot.Version{}, errors.Join(ErrParseVersion, err)
	}

	return version, nil
}

/* ---------------------------- Function: Resolve --------------------------- */

// Tries to locate a pin file in the current directory or any parent directories.
func Resolve(path string) (string, error) {
	// Check if the specified path (or any ancestors) has a pin
	for path != "/" {
		// Don't overwrite 'path' or you'll go into an infinite loop due to
		// 'Clean()' appending filenames you're removing below.
		pin, err := Clean(path)
		if err != nil {
			return "", err
		}

		info, err := os.Stat(pin)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return "", err
			}

			path = filepath.Dir(path)

			continue
		}

		// Validate that the file is a regular file; this catches cases where
		// there's a directory named after 'pinFilename'.
		if info.Mode().IsRegular() {
			return pin, nil
		}
	}

	return "", fs.ErrNotExist
}

/* ---------------------------- Function: Remove ---------------------------- */

// Deletes the specified pin file if it exists.
func Remove(path string) error {
	p, err := Clean(path)
	if err != nil {
		return err
	}

	if err := os.Remove(p); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}

/* ----------------------------- Function: Write ---------------------------- */

// Writes a 'Version' to the specified pin file path.
func Write(version godot.Version, path string) error {
	path, err := Clean(path)
	if err != nil {
		return err
	}

	// Make the parent directories if needed.
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(version.String()), os.ModePerm)
}
