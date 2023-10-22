package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const envStore = "GDENV_HOME"
const storeName = "gdenv"

var (
	ErrIllegalPath   = errors.New("illegal store path")
	ErrInvalidPath   = errors.New("invalid file path")
	ErrMissingEnvVar = errors.New("missing environment variable")
)

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
/*                               Function: Path                               */
/* -------------------------------------------------------------------------- */

// Returns the user-configured path to the 'gdenv' store.
func Path() (string, error) {
	path := os.Getenv(envStore)
	if path == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingEnvVar, envStore)
	}

	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("%w; expected absolute path: %s", ErrInvalidPath, path)
	}

	if base := filepath.Base(path); base != storeName && base != "."+storeName {
		return "", fmt.Errorf("%w: '%s'", ErrIllegalPath, path)
	}

	return filepath.Clean(path), nil
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
