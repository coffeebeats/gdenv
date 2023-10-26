package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/download"
	"github.com/coffeebeats/gdenv/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Downloads and caches a platform-specific version of Godot.
func Executable(ctx context.Context, storePath string, ex executable.Executable) error {
	m, err := mirror.Choose(ctx, ex.Version(), ex.Platform())
	if err != nil {
		return err
	}

	log.Infof("downloading from mirror: %s", strings.TrimPrefix(fmt.Sprintf("%T", m), "mirror."))

	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	log.Debugf("using temporary directory: %s", tmp)

	localExArchive, err := download.ExecutableWithChecksumValidation(ctx, m, ex, tmp)
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
	// TODO: Make this not rely on this (arbitrary) platform. It would be better
	// if 'mirror.checkIfExists' could correctly determine existence of an
	// arbitrary artifact. For now, select a platform that will certainly exist.
	p := platform.Platform{Arch: platform.Amd64, OS: platform.Windows}

	m, err := mirror.Choose(ctx, src.Version(), p)
	if err != nil {
		return err
	}

	log.Infof("downloading from mirror: %s", strings.TrimPrefix(fmt.Sprintf("%T", m), "mirror."))

	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	log.Debugf("using temporary directory: %s", tmp)

	localSourceArchive, err := download.SourceWithChecksumValidation(ctx, m, src.Version(), tmp)
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
