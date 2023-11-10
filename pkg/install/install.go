package install

import (
	"context"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/pkg/download"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Downloads and caches a platform-specific version of Godot.
func Executable(ctx context.Context, storePath string, ex executable.Executable) error {
	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	log.Debugf("using temporary directory: %s", tmp)

	localExArchive, err := download.ExecutableWithChecksumValidation(ctx, ex, tmp)
	if err != nil {
		return err
	}

	log.Info("installing executable in gdenv store")

	if err := archive.Extract[executable.Archive](ctx, localExArchive, tmp); err != nil {
		return err
	}

	if err := os.Remove(localExArchive.Path); err != nil {
		return err
	}

	log.Debug("extracted executable archive")

	entries, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}

	artifacts := make([]artifact.Local[artifact.Artifact], 0, len(entries))
	for _, entry := range entries {
		artifacts = append(artifacts, artifact.Local[artifact.Artifact]{
			Artifact: ex,
			Path:     filepath.Join(tmp, entry.Name()),
		})
	}

	return store.Add(storePath, artifacts...)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Source                              */
/* -------------------------------------------------------------------------- */

// Downloads and caches a specific version of Godot's source code.
func Source(ctx context.Context, storePath string, src source.Source) error {
	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	log.Debugf("using temporary directory: %s", tmp)

	localSourceArchive, err := download.SourceWithChecksumValidation(ctx, src.Version(), tmp)
	if err != nil {
		return err
	}

	log.Info("installing source in gdenv store")

	return store.Add(
		storePath,
		artifact.Local[artifact.Artifact]{
			Artifact: localSourceArchive.Artifact,
			Path:     localSourceArchive.Path,
		},
	)
}
