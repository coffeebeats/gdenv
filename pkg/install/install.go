package install

import (
	"context"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/pkg/download"
	"github.com/coffeebeats/gdenv/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Downloads and caches a platform-specific version of Godot.
func Executable(ctx context.Context, storePath string, ex executable.Executable) error {
	m, err := mirror.Choose(ctx, ex.Version())
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

// Downloads and caches a platform-specific version of Godot.
func Source(ctx context.Context, storePath string, s source.Source) error {
	m, err := mirror.Choose(ctx, s.Version())
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
