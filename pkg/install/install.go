package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/pkg/download"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Downloads and caches a platform-specific version of Godot.
func Executable( //nolint:funlen
	ctx context.Context,
	storePath string,
	ex executable.Executable,
	force bool,
) error {
	p, v := ex.Platform(), ex.Version()

	ok, err := store.Has(storePath, ex)
	if err != nil {
		return err
	}

	if ok && !force {
		log.Info("skipping installation; version already found")

		return nil
	}

	platformLabel, err := platform.Format(p, v)
	if err != nil {
		return fmt.Errorf("%w: %w", platform.ErrUnrecognizedPlatform, err)
	}

	log.Infof("installing version: %s (%s)", v, platformLabel)

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

	log.Info("adding executable to gdenv store")

	if err := archive.Extract[executable.Archive](ctx, localExArchive, tmp); err != nil {
		return err
	}

	if err := os.Remove(localExArchive.Path); err != nil {
		return err
	}

	log.Debug("successfully extracted executable archive")

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

	if err := store.Add(ctx, storePath, artifacts...); err != nil {
		return err
	}

	log.Infof("successfully installed version: %s (%s,%s)", v, p.OS, p.Arch)

	return nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: Source                              */
/* -------------------------------------------------------------------------- */

// Downloads and caches a specific version of Godot's source code.
func Source(ctx context.Context, storePath string, v version.Version, force bool) error {
	// Ensure the store exists.
	if err := store.Touch(storePath); err != nil {
		return err
	}

	// Define the target 'Source'.
	src := source.New(v)

	ok, err := store.Has(storePath, src)
	if err != nil {
		return err
	}

	if ok && !force {
		log.Info("skipping installation; version already found")

		return nil
	}

	log.Infof("installing version: %s", v)

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

	log.Debug("installing source in gdenv store")

	if err := store.Add(
		ctx,
		storePath,
		artifact.Local[artifact.Artifact]{
			Artifact: localSourceArchive.Artifact,
			Path:     localSourceArchive.Path,
		},
	); err != nil {
		return err
	}

	log.Infof("successfully installed version: %s", src.Version())

	return nil
}
