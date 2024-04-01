package store

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
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
type LocalSrc = artifact.Local[source.Archive]

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
		if err := os.MkdirAll(pathArtifactDir, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		path := filepath.Join(pathArtifactDir, filepath.Base(local.Path))
		if err := osutil.CopyFile(context.TODO(), local.Path, path); err != nil {
			return err
		}

		log.Debugf("added file to store: %s", path)
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
/*                            Function: Executables                           */
/* -------------------------------------------------------------------------- */

// Executables returns the list of installed Godot executables.
func Executables(ctx context.Context, storePath string) ([]LocalEx, error) {
	if storePath == "" {
		return nil, ErrMissingStore
	}

	out := make([]LocalEx, 0)

	entries, err := os.ReadDir(filepath.Join(storePath, storeDirEx))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		return nil, nil
	}

	for _, versionDir := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		v, err := version.Parse(versionDir.Name())
		if err != nil {
			continue
		}

		ee, err := collectExecutablesForVersion(ctx, storePath, v)
		if err != nil {
			return nil, err
		}

		out = append(out, ee...)
	}

	return out, nil
}

// Returns all cached executables for the specified version.
func collectExecutablesForVersion(
	ctx context.Context,
	storePath string,
	v version.Version,
) ([]LocalEx, error) {
	out := make([]LocalEx, 0)

	entries, err := os.ReadDir(filepath.Join(storePath, storeDirEx, v.String()))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		return nil, nil
	}

	for _, platformDir := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		p, err := platform.Parse(platformDir.Name())
		if err != nil {
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

		path, err := artifactPath(storePath, ex)
		if err != nil {
			return nil, err
		}

		out = append(out, LocalEx{Artifact: ex, Path: path})
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

	// For executables, the entire platform directory should be removed. This is
	// because multiple files can be installed alongside the executable itself.
	if _, ok := a.(executable.Executable); ok {
		path = filepath.Dir(path)
	}

	// Remove the specific executable from the store.
	if err := os.RemoveAll(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	log.Debugf("removed directory from store: %s", path)

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
/*                              Function: Sources                             */
/* -------------------------------------------------------------------------- */

// Sources returns the list of installed Godot source code versions.
func Sources(ctx context.Context, storePath string) ([]LocalSrc, error) {
	if storePath == "" {
		return nil, ErrMissingStore
	}

	out := make([]LocalSrc, 0)

	entries, err := os.ReadDir(filepath.Join(storePath, storeDirSrc))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		return nil, nil
	}

	for _, versionDir := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		v, err := version.Parse(versionDir.Name())
		if err != nil {
			continue
		}

		src := source.New(v)

		ok, err := Has(storePath, src)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		path, err := artifactPath(storePath, src)
		if err != nil {
			return nil, err
		}

		out = append(
			out,
			LocalSrc{Artifact: source.Archive{Inner: src}, Path: path},
		)
	}

	return out, nil
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
	if err := os.MkdirAll(storePath, osutil.ModeUserRWXGroupRX); err != nil {
		return err
	}

	// Create the required subdirectories, if needed.
	for _, d := range []string{storeDirBin, storeDirSrc, storeDirEx} {
		path := filepath.Join(storePath, d)
		if err := os.MkdirAll(path, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}
	}

	// Create the required files, if needed.
	for _, f := range []string{storeFileLayout} {
		path := filepath.Join(storePath, f)
		if err := os.WriteFile(path, nil, osutil.ModeUserRW); err != nil {
			return err
		}
	}

	return nil
}
