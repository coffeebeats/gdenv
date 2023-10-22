package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/internal/osutil"
)

const (
	modeStoreDir    = 0755 // rwxr-xr-x
	modeStoreLayout = 0644 // rw-r--r--

	storeDirBin     = "bin"
	storeDirSrc     = "src"
	storeDirEx      = "editor"
	storeFileLayout = "layout.v0" // simplify migrating in the future
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrMissingStore        = errors.New("missing store")
	ErrUnexpectedLayout    = errors.New("unexpected layout")
	ErrUnsupportedArtifact = errors.New("unsupported artifact")
)

type LocalEx = artifact.Local[executable.Executable]

/* -------------------------------------------------------------------------- */
/*                                Function: Add                               */
/* -------------------------------------------------------------------------- */

// Add caches the specified locally-available artifacts in the store.
func Add(storePath string, localArtifacts ...artifact.Local[artifact.Artifact]) error {
	if storePath == "" {
		return ErrMissingStore
	}

	// Verify that the files-to-add exist.
	for _, a := range localArtifacts {
		if _, err := os.Stat(a.Path); err != nil {
			return err
		}
	}

	// Add the specified artifacts to the store.
	for _, local := range localArtifacts {
		// Determine the directory to place the files under.
		pathArtifact, err := artifactPath(storePath, local.Artifact)
		if err != nil {
			return err
		}

		pathArtifactDir := filepath.Dir(pathArtifact)

		// Create the required directories, if needed.
		if err := os.MkdirAll(pathArtifactDir, modeStoreDir); err != nil {
			return err
		}

		path := filepath.Join(pathArtifactDir, filepath.Base(local.Path))
		if err := osutil.ForceRename(local.Path, path); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Function: Clear                              */
/* -------------------------------------------------------------------------- */

// Removes all cached artifacts in the store.
func Clear(storePath string) error {
	if storePath == "" {
		return ErrMissingStore
	}

	// Clear the entire source cache directory.
	if err := os.RemoveAll(filepath.Join(storePath, storeDirSrc)); err != nil {
		return err
	}

	// Clear the entire executable cache directory.
	if err := os.RemoveAll(filepath.Join(storePath, storeDirEx)); err != nil {
		return err
	}

	// Remake the deleted directories.
	return Touch(storePath)
}

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Returns the full path (starting with the store path) to the *executable* file
// in the store.
//
// NOTE: This does *not* mean the executable exists.
func Executable(storePath string, ex executable.Executable) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	pathExecutableDir, err := executableDir(storePath, ex)
	if err != nil {
		return "", err
	}

	return filepath.Join(pathExecutableDir, ex.Path()), nil
}

/* -------------------------------------------------------------------------- */
/*                            Function: Executables                           */
/* -------------------------------------------------------------------------- */

// Executables returns the list of installed Godot executables.
func Executables(ctx context.Context, storePath string) ([]LocalEx, error) { //nolint:cyclop
	if storePath == "" {
		return nil, ErrMissingStore
	}

	out := make([]LocalEx, 0)

	entries, err := os.ReadDir(filepath.Join(storePath, storeDirEx))
	if err != nil {
		return nil, err
	}

	for _, versionDir := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		v, err := version.Parse(versionDir.Name())
		if err != nil {
			log.Println("Skipping directory", versionDir.Name())
			continue
		}

		entries, err := os.ReadDir(filepath.Join(storePath, storeDirEx, versionDir.Name()))
		if err != nil {
			return nil, err
		}

		for _, platformDir := range entries {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			p, err := platform.Parse(platformDir.Name())
			if err != nil {
				log.Println("Skipping directory", versionDir.Name())
				continue
			}

			ex := executable.New(v, p)

			ok, err := Has(storePath, ex)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}

			path, err := Executable(storePath, ex)
			if err != nil {
				return nil, err
			}

			out = append(out, LocalEx{Artifact: ex, Path: path})
		}
	}

	return out, nil
}

/* -------------------------------------------------------------------------- */
/*                                Function: Has                               */
/* -------------------------------------------------------------------------- */

// Return whether the store has the specified version cached.
func Has(storePath string, a artifact.Artifact) (bool, error) {
	if storePath == "" {
		return false, ErrMissingStore
	}

	path, err := artifactPath(storePath, a)
	if err != nil {
		if !errors.Is(err, ErrUnsupportedArtifact) {
			return false, err
		}

		return false, nil
	}

	_, err = os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: Remove                              */
/* -------------------------------------------------------------------------- */

// Removes the specified version from the store.
func Remove(storePath string, a artifact.Artifact) error {
	if storePath == "" {
		return ErrMissingStore
	}

	path, err := artifactPath(storePath, a)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, ErrUnsupportedArtifact) {
			return nil
		}

		return err
	}

	// Remove the specific executable from the store.
	if err := os.Remove(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return removeUnusedCacheDirectories(storePath, path)
}

// A utility method which cleans up unused directories from the specified path
// up to the store's cache directories.
func removeUnusedCacheDirectories(storePath, path string) error {
	for {
		path = filepath.Dir(path)

		// Add a safeguard to not escape the store cache directories.
		if path == filepath.Join(storePath, storeDirEx) || path == filepath.Join(storePath, storeDirSrc) {
			return nil
		}

		files, err := os.ReadDir(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if len(files) > 0 {
			return nil
		}

		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                              Function: Source                              */
/* -------------------------------------------------------------------------- */

// Returns the full path (starting with the store path) to the Godot source
// directory in the store.
//
// NOTE: This does *not* mean the source folder exists.
func Source(storePath string, src source.Source) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	return artifactPath(storePath, src)
}

/* -------------------------------------------------------------------------- */
/*                               Function: Touch                              */
/* -------------------------------------------------------------------------- */

// Touch ensures a store is initialized at the specified path; no effect if it
// exists already.
func Touch(storePath string) error {
	if storePath == "" {
		return ErrMissingStore
	}

	// Create the 'Store' directory, if needed.
	if err := os.MkdirAll(storePath, modeStoreDir); err != nil {
		return err
	}

	// Create the required subdirectories, if needed.
	for _, d := range []string{storeDirBin, storeDirSrc, storeDirEx} {
		path := filepath.Join(storePath, d)
		if err := os.MkdirAll(path, modeStoreDir); err != nil {
			return err
		}
	}

	// Create the required files, if needed.
	for _, f := range []string{storeFileLayout} {
		path := filepath.Join(storePath, f)
		if err := os.WriteFile(path, nil, modeStoreLayout); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: artifactPath                           */
/* -------------------------------------------------------------------------- */

// artifactPath returns the path (starting with the store path) to the artifact
// cached in the store.
//
// NOTE: This does *not* mean the artifact exists.
func artifactPath(storePath string, a artifact.Artifact) (string, error) {
	switch a := a.(type) {
	case executable.Executable:
		path, err := executableDir(storePath, a)
		if err != nil {
			return "", err
		}

		pathExParts := strings.Split(a.Path(), string(os.PathSeparator))
		if len(pathExParts) == 0 {
			return "", fmt.Errorf("%w: missing executable path: '%s'", ErrInvalidInput, a.Path())
		}

		return filepath.Join(path, pathExParts[0]), nil
	case source.Source:
		pathSourceDir, err := sourceDir(storePath, a)
		if err != nil {
			return "", err
		}

		return filepath.Join(pathSourceDir, a.Name()), nil
	}

	return "", fmt.Errorf("%w: %T", ErrUnsupportedArtifact, a)
}

/* ------------------------- Function: executableDir ------------------------ */

func executableDir(storePath string, ex executable.Executable) (string, error) {
	if err := version.Validate(ex.Version()); err != nil {
		return "", err
	}

	platformLabel, err := platform.Format(ex.Platform(), ex.Version())
	if err != nil {
		return "", fmt.Errorf("%w: missing platform: %w", ErrInvalidInput, err)
	}

	path := filepath.Join(
		storePath,
		storeDirEx,
		ex.Version().String(),
		platformLabel,
	)

	return path, nil
}

/* --------------------------- Function: sourceDir -------------------------- */

func sourceDir(storePath string, src source.Source) (string, error) {
	if err := version.Validate(src.Version()); err != nil {
		return "", err
	}

	path := filepath.Join(storePath, storeDirSrc, src.Version().String())

	return path, nil
}
