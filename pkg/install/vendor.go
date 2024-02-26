package install

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/store"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                              Function: Vendor                              */
/* -------------------------------------------------------------------------- */

// Extracts the cached source code folder into the specified 'out' path.
//
// NOTE: This function will fail if the source code does not exist in the store.
func Vendor(ctx context.Context, v version.Version, storePath, out string) error {
	src := source.Archive{Inner: source.New(v)}

	srcPath, err := store.Source(storePath, src.Inner)
	if err != nil {
		return err
	}

	if out == "" {
		return fmt.Errorf("%w: vendor directory", ErrMissingInput)
	}

	out = filepath.Clean(out)

	// Improve log clarity by prefixing a relative path with './'.
	if !filepath.IsAbs(out) && !strings.HasPrefix(out, "..") && !strings.HasPrefix(out, "./") {
		out = "./" + out
	}

	info, err := os.Stat(out)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err := os.MkdirAll(out, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("%w: %s", fs.ErrExist, out)
	}

	localSrcArchive := artifact.Local[source.Archive]{Artifact: src, Path: srcPath}
	if err := archive.Extract(ctx, localSrcArchive, out); err != nil {
		return err
	}

	log.Infof("successfully vendored version %s: %s", v, out)

	return nil
}
