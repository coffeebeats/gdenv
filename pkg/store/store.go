package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot"
)

var (
	ErrIOFailed         = errors.New("store: IO failed")
	ErrDirNotFound      = errors.New("store: directory not found")
	ErrFileNotFound     = errors.New("store: file not found")
	ErrUnexpectedLayout = errors.New("store: unexpected layout")

	godotDir = "godot"
)

/* ------------------------------ Function: Add ----------------------------- */

// Move the specified file into the store for the specified version.
func Add(p, t string, v godot.Version) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if !Exists(p) {
		return fmt.Errorf("%w: %s", ErrDirNotFound, p)
	}

	name, err := godot.Executable(v)
	if err != nil {
		return err
	}

	execPath := filepath.Join(p, godotDir, v.String(), name)

	if _, err := os.Stat(t); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Join(ErrIOFailed, err)
		}
	}

	// Create any parent directories required.
	if err := os.MkdirAll(filepath.Base(execPath), os.ModeDir); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	// Move the file into the store.
	if err := os.Rename(t, execPath); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}

/* ----------------------------- Function: Find ----------------------------- */

// Returns the path to the specified version of Godot; if it does not exist,
// an error is returned.
func Find(p string, v godot.Version) (string, error) {
	p, err := Clean(p)
	if err != nil {
		return "", err
	}

	if !Exists(p) {
		return "", fmt.Errorf("%w: %s", ErrDirNotFound, p)
	}

	name, err := godot.Executable(v)
	if err != nil {
		return "", err
	}

	execPath := filepath.Join(p, godotDir, v.String(), name)

	info, err := os.Stat(execPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", errors.Join(ErrFileNotFound, err)
		}

		return "", errors.Join(ErrIOFailed, err)
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%w: %s", ErrUnexpectedLayout, execPath)
	}

	return execPath, nil
}

/* ------------------------------ Function: Has ----------------------------- */

// Returns whether or not the store contains the specified version of Godot.
func Has(p string, v godot.Version) (bool, error) {
	_, err := Find(p, v)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

/* ----------------------------- Function: Init ----------------------------- */

// Initializes a store at the specified path; no effect if it exists already.
func Init(p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	// Create the required subdirectories.
	if err := os.MkdirAll(p, os.ModeDir); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	for _, d := range []string{"bin", "godot"} {
		if err := os.MkdirAll(filepath.Join(p, d), os.ModeDir); err != nil {
			return errors.Join(ErrIOFailed, err)
		}
	}

	return nil
}

/* ----------------------------- Function: List ----------------------------- */

// Returns a list of cached versions of Godot.
func List(p string, v godot.Version) ([]godot.Version, error) {
	p, err := Clean(p)
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(filepath.Join(p, godotDir))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIOFailed, err)
	}

	out := make([]godot.Version, len(dirs))

	for i, d := range dirs {
		v, err := godot.ParseVersion(d.Name())
		if err != nil {
			return nil, err
		}

		// Check that the executable for the current platform exists.
		if _, err := Find(p, v); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}

			continue
		}

		out[i] = v
	}

	return out, nil
}

/* ---------------------------- Function: Remove ---------------------------- */

// Removes the specified version from the store.
func Remove(p string, v godot.Version) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if !Exists(p) {
		return fmt.Errorf("%w: %s", ErrDirNotFound, p)
	}

	execPath, err := Find(p, v)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	// Remove the specific executable from the store.
	if err := os.Remove(execPath); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	// Check if the parent directory is empty. If it is, remove it.
	d := filepath.Base(execPath)

	files, err := os.ReadDir(d)
	if err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	if len(files) > 0 {
		return nil
	}

	if err := os.RemoveAll(filepath.Base(execPath)); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}
