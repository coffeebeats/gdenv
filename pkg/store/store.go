package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	modeStoreDir    = 0755 // rwxr-xr-x
	modeStoreLayout = 0644 // rw-r--r--

	storeDirBin     = "bin"
	storeDirGodot   = "godot"
	storeFileLayout = "layout.v0" // simplify migrating in the future
)

var (
	ErrInvalidSpecification = errors.New("invalid specification")
	ErrMissingStore         = errors.New("missing store")
	ErrUnexpectedLayout     = errors.New("unexpected layout")
)

/* ----------------------------- Function: Init ----------------------------- */

// Initializes a store at the specified path; no effect if it exists already.
func Init(path string) error {
	path, err := Clean(path)
	if err != nil {
		return err
	}

	// Create the 'Store' directory, if needed.
	if err := os.MkdirAll(path, modeStoreDir); err != nil {
		return err
	}

	// Create the required subdirectories, if needed.
	for _, d := range []string{storeDirBin, storeDirGodot} {
		if err := os.MkdirAll(filepath.Join(path, d), modeStoreDir); err != nil {
			return err
		}
	}

	// Create the required files, if needed.
	for _, f := range []string{storeFileLayout} {
		if err := os.WriteFile(filepath.Join(path, f), nil, modeStoreLayout); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------- Function: InitAtPath -------------------------- */

// A convenience method which initializes a store at the path specified by
// the 'envVarStore' environment variable.
func InitAtPath() (string, error) {
	storePath, err := Path()
	if err != nil {
		return "", err
	}

	if err := Init(storePath); err != nil {
		return "", err
	}

	return storePath, nil
}

/* ------------------------------ Function: Add ----------------------------- */

// Move the specified file into the store for the specified version.
func Add(store, file string, ex executable.Executable) error {
	store, err := Clean(store)
	if err != nil {
		return err
	}

	if !Exists(store) {
		return fmt.Errorf("%w: '%s'", ErrMissingStore, store)
	}

	tool, err := ToolPath(store, ex)
	if err != nil {
		return err
	}

	// Create the required directories, if needed.
	if err := os.MkdirAll(filepath.Dir(tool), modeStoreDir); err != nil {
		return err
	}

	return os.Rename(file, tool)
}

/* ------------------------------ Function: Has ----------------------------- */

// Return whether the store has the specified version cached.
func Has(store string, ex executable.Executable) bool {
	store, err := Clean(store)
	if err != nil {
		return false
	}

	if !Exists(store) {
		return false
	}

	tool, err := ToolPath(store, ex)
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
func Remove(store string, ex executable.Executable) error {
	store, err := Clean(store)
	if err != nil {
		return err
	}

	if !Exists(store) {
		return fmt.Errorf("%w: '%s'", ErrMissingStore, store)
	}

	if !Has(store, ex) {
		return nil
	}

	tool, err := ToolPath(store, ex)
	if err != nil {
		return err
	}

	// Remove the specific executable from the store.
	if err := os.Remove(tool); err != nil {
		return err
	}

	// Check if the parent directory is empty. If it is, remove it.
	toolParent := filepath.Dir(tool)

	files, err := os.ReadDir(toolParent)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		return nil
	}

	return os.RemoveAll(toolParent)
}

/* --------------------------- Function: ToolPath --------------------------- */

// Returns the full path to the tool in the store.
//
// NOTE: This does *not* mean the tool exists.
func ToolPath(store string, ex executable.Executable) (string, error) {
	store, err := Clean(store)
	if err != nil {
		return "", err
	}

	return filepath.Join(store, storeDirGodot, ex.Version().String(), ex.Name()), nil
}

/* --------------------------- Function: Versions --------------------------- */

// Returns a list of cached Godot executables.
func Executables(ctx context.Context, store string) ([]executable.Executable, error) {
	store, err := Clean(store)
	if err != nil {
		return nil, err
	}

	if !Exists(store) {
		return nil, fmt.Errorf("%w: '%s'", ErrMissingStore, store)
	}

	pathCache := filepath.Join(store, storeDirGodot)

	versions, err := os.ReadDir(pathCache)
	if err != nil {
		// NOTE: Some versions *may* have been found in this case, as 'ReadDir'
		// returns what it can, but it's safer to just fail here entirely.
		return nil, fmt.Errorf("%w: '%s'", err, pathCache)
	}

	out := make([]executable.Executable, 0)

	for _, dirVersion := range versions {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}

		v, err := version.Parse(dirVersion.Name())
		if err != nil {
			return nil, errors.Join(ErrInvalidSpecification, err)
		}

		pathVersion := filepath.Join(pathCache, dirVersion.Name())

		executables, err := collectExecutables(pathVersion)
		if err != nil {
			return out, err
		}

		for _, ex := range executables {
			if ctx.Err() != nil {
				return out, ctx.Err()
			}

			// Check that the executable for the current platform exists. Much
			// of this call is redundant with checks above, but 'Has' may also
			// check things like whether the file is indeed a file, etc.).
			if !Has(store, ex) {
				return out, fmt.Errorf("%w: '%s'", ErrUnexpectedLayout, ex)
			}

			// Validate that executables are found under the correct directory.
			if ex.Version() != v {
				return out, fmt.Errorf("%w: '%s'", ErrUnexpectedLayout, ex)
			}
		}

		out = append(out, executables...)
	}

	return out, nil
}

/* ---------------------- Function: collectExecutables ---------------------- */

// Collects the set of 'Executable' files found under the specified directory.
func collectExecutables(path string) ([]executable.Executable, error) {
	out := make([]executable.Executable, 0)

	executables, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("%w: '%s'", err, path)
	}

	for _, file := range executables {
		ex, err := executable.Parse(file.Name())
		if err != nil {
			return nil, errors.Join(ErrInvalidSpecification, err)
		}

		out = append(out, ex)
	}

	return out, nil
}
