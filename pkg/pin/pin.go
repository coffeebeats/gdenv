package pin

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/internal/pathutil"
)

const modePinFile = 0664 // rw-rw-r--

var (
	ErrParseVersion   = errors.New("failed to parse version")
	ErrUnexpectedFile = errors.New("unexpected file")
)

/* ----------------------------- Function: Read ----------------------------- */

// Parses a 'Version' from the specified pin file.
func Read(path string) (version.Version, error) {
	path, err := Clean(path)
	if err != nil {
		return version.Version{}, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return version.Version{}, err
	}

	v, err := version.Parse(string(bytes))
	if err != nil {
		return version.Version{}, errors.Join(ErrParseVersion, err)
	}

	return v, nil
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
func Write(v version.Version, path string) error {
	path, err := Clean(path)
	if err != nil {
		return err
	}

	// Determine the permissions of the nearest ancestor directory.
	mode, err := pathutil.AncestorMode(path)
	if err != nil {
		return fmt.Errorf("cannot determine permissions: %w", err)
	}

	// Make the parent directories if needed.
	if err := os.MkdirAll(filepath.Dir(path), mode); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(v.String()), modePinFile)
}
