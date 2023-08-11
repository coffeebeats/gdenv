package pin

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot"
)

var (
	ErrIOFailed       = errors.New("pin: IO failed")
	ErrFileNotFound   = errors.New("pin: file not found")
	ErrUnexpectedFile = errors.New("pin: unexpected file")
)

/* ----------------------------- Function: Read ----------------------------- */

// Parses a 'Version' from the specified pin file.
func Read(p string) (godot.Version, error) {
	p, err := Clean(p)
	if err != nil {
		return godot.Version{}, err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return godot.Version{}, errors.Join(ErrIOFailed, err)
	}

	return godot.ParseVersion(string(b))
}

/* ---------------------------- Function: Resolve --------------------------- */

// Tries to locate a pin file in the current directory or any parent directories.
func Resolve(p string) (string, error) {
	var path = p

	// Check if the specified path (or any ancestors) has a pin
	for len(path) > 0 {
		// Don't overwrite 'path' or you'll go into an infinite loop due to
		// 'Pin.Path()' appending filenames you're removing below.
		p, err := Clean(path)
		if err != nil {
			return "", err
		}

		info, err := os.Stat(p)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return "", errors.Join(ErrIOFailed, err)
			}

			d, _ := filepath.Split(path)
			path = d

			continue
		}

		// Validate that the file is a regular file; this catches cases where
		// there's a directory named after 'pinFilename'.
		if info.Mode().IsRegular() {
			return p, nil
		}
	}

	return "", ErrFileNotFound
}

/* ----------------------------- Function: Write ---------------------------- */

// Deletes the specified pin file.
func Remove(p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if err := os.Remove(p); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Join(ErrIOFailed, err)
		}
	}

	return nil
}

/* ----------------------------- Function: Write ---------------------------- */

// Writes a 'Version' to the specified pin file path.
func Write(v godot.Version, p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if err := os.WriteFile(p, []byte(v.String()), 0); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}
