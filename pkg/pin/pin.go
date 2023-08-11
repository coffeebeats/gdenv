package pin

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot"
)

var (
	ErrIOFailed       = errors.New("pin: IO failed")
	ErrFileNotFound   = errors.New("pin: file not found")
	ErrParseVersion   = errors.New("pin: failed to parse version")
	ErrUnexpectedFile = errors.New("pin: unexpected file")
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
		if errors.Is(err, fs.ErrNotExist) {
			return godot.Version{}, errors.Join(ErrFileNotFound, err)
		}

		return godot.Version{}, errors.Join(ErrIOFailed, err)
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
				return "", errors.Join(ErrIOFailed, err)
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

	return "", ErrFileNotFound
}

/* ----------------------------- Function: Write ---------------------------- */

// Deletes the specified pin file if it exists.
func Remove(path string) error {
	p, err := Clean(path)
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
func Write(version godot.Version, path string) error {
	path, err := Clean(path)
	if err != nil {
		return err
	}

	log.Print(filepath.Dir(path))

	// Make the parent directories if needed.
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	contents := version.Canonical().String()
	if err := os.WriteFile(path, []byte(contents), os.ModePerm); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}
