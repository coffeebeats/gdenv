package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot"
)

const storeDirBin = "bin"
const storeDirGodot = "godot"
const storeFileLayout = "layout.v1" // simplify migrating in the future

var (
	ErrDirNotFound      = errors.New("store: directory not found")
	ErrFileNotFound     = errors.New("store: file not found")
	ErrIOFailed         = errors.New("store: I/O failed")
	ErrInvalidVersion   = errors.New("store: invalid version")
	ErrMissingStore     = errors.New("store: missing store")
	ErrUnexpectedLayout = errors.New("store: unexpected layout")
)

/* ----------------------------- Function: Init ----------------------------- */

// Initializes a store at the specified path; no effect if it exists already.
func Init(path string) error {
	path, err := Clean(path)
	if err != nil {
		return err
	}

	// Create the 'Store' directory, if needed.
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	// Create the required subdirectories, if needed.
	for _, d := range []string{storeDirBin, storeDirGodot} {
		if err := os.MkdirAll(filepath.Join(path, d), os.ModePerm); err != nil {
			return errors.Join(ErrIOFailed, err)
		}
	}

	// Create the required files, if needed.
	for _, f := range []string{storeFileLayout} {
		if err := os.WriteFile(filepath.Join(path, f), nil, os.ModePerm); err != nil {
			return errors.Join(ErrIOFailed, err)
		}
	}

	return nil
}

/* ------------------------------ Function: Add ----------------------------- */

// Move the specified file into the store for the specified version.
func Add(store, file string, version godot.Version) error {
	store, err := Clean(store)
	if err != nil {
		return err
	}

	if !Exists(store) {
		return fmt.Errorf("%w: %s", ErrMissingStore, store)
	}

	tool, err := ToolPath(store, version)
	if err != nil {
		return err
	}

	// Create the required directories, if needed.
	if err := os.MkdirAll(filepath.Dir(tool), os.ModePerm); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	if err := os.Rename(file, tool); err != nil {
		log.Print(err)
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}

/* ------------------------------ Function: Has ----------------------------- */

// Return whether the store has the specified version cached.
func Has(store string, version godot.Version) bool {
	store, err := Clean(store)
	if err != nil {
		return false
	}

	if !Exists(store) {
		return false
	}

	tool, err := ToolPath(store, version)
	if err != nil {
		return false
	}

	info, err := os.Stat(tool)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

/* ---------------------------- Function: Remove ---------------------------- */

// Removes the specified version from the store.
func Remove(store string, version godot.Version) error {
	store, err := Clean(store)
	if err != nil {
		return err
	}

	if !Exists(store) {
		return fmt.Errorf("%w: %s", ErrMissingStore, store)
	}

	if !Has(store, version) {
		return nil
	}

	tool, err := ToolPath(store, version)
	if err != nil {
		return err
	}

	// Remove the specific executable from the store.
	if err := os.Remove(tool); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	// Check if the parent directory is empty. If it is, remove it.
	toolParent := filepath.Dir(tool)

	files, err := os.ReadDir(toolParent)
	if err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	if len(files) > 0 {
		return nil
	}

	if err := os.RemoveAll(toolParent); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}

/* --------------------------- Function: ToolPath --------------------------- */

// Returns the full path to the tool in the store.
//
// NOTE: This does *not* mean the tool exists.
func ToolPath(store string, version godot.Version) (string, error) {
	version = version.Canonical()

	store, err := Clean(store)
	if err != nil {
		return "", err
	}

	name, err := godot.ExecutableName(version)
	if err != nil {
		return "", errors.Join(ErrInvalidVersion, err)
	}

	return filepath.Join(store, storeDirGodot, version.String(), name), nil
}

/* --------------------------- Function: Versions --------------------------- */

// Returns a list of cached versions of Godot.
func Versions(store string) ([]godot.Version, error) {
	store, err := Clean(store)
	if err != nil {
		return nil, err
	}

	if !Exists(store) {
		return nil, fmt.Errorf("%w: %s", ErrMissingStore, store)
	}

	cache := filepath.Join(store, storeDirGodot)

	dirs, err := os.ReadDir(cache)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDirNotFound, cache)
	}

	out := make([]godot.Version, 0)

	for _, d := range dirs {
		version, err := godot.ParseVersion(d.Name())
		if err != nil {
			return nil, errors.Join(ErrInvalidVersion, err)
		}

		// Check that the executable for the current platform exists.
		if !Has(store, version) {
			continue
		}

		out = append(out, version)
	}

	return out, nil
}
