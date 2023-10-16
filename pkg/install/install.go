package install

import (
	"context"
	"os"

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

	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	localExArchive, err := download.ExecutableWithChecksumValidation(ctx, m, ex, tmp)
	if err != nil {
		return err
	}

	if err := archive.Extract[executable.Archive](ctx, localExArchive, tmp); err != nil {
		return err
	}

	if err := os.Remove(localExArchive.Path); err != nil {
		return err
	}

	return store.AddDirectory(storePath, ex.Version(), tmp)
}

/* -------------------------------------------------------------------------- */
/*                              Function: Source                              */
/* -------------------------------------------------------------------------- */

// Downloads and caches a specific version of Godot's source code.
func Source(ctx context.Context, storePath string, s source.Source) error {
	// TODO: Make this not rely on this (arbitrary) platform. It would be better
	// if 'checkIfExists' could correctly determine existence of an arbitrary
	// artifact. For now, select a platform that's definitely going to exist.
	p := platform.Platform{Arch: platform.Amd64, OS: platform.Windows}

	m, err := mirror.Choose(ctx, s.Version(), p)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "gdenv-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	localSourceArchive, err := download.SourceWithChecksumValidation(ctx, m, s.Version(), tmp)
	if err != nil {
		return err
	}

	if err := archive.Extract[source.Archive](ctx, localSourceArchive, tmp); err != nil {
		return err
	}

	if err := os.Remove(localSourceArchive.Path); err != nil {
		return err
	}

	return store.AddDirectory(storePath, s.Version(), tmp)
}
