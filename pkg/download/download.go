package download

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/mirror"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

type progressKey[T artifact.Artifact] struct{}

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithProgress creates a sub-context with an associated progress reporter. The
// result can be passed to download functions in this package to get updates on
// the download progress for that specific artifact.
func WithProgress[T artifact.Artifact](
	ctx context.Context,
	p *progress.Progress,
) context.Context {
	return context.WithValue(ctx, progressKey[T]{}, p)
}

/* -------------------------------------------------------------------------- */
/*                             Function: Download                             */
/* -------------------------------------------------------------------------- */

// Download uses the provided mirror to download the specified artifact and
// returns an 'artifact.Local' wrapper pointing to it.
func Download[T artifact.Artifact](
	ctx context.Context,
	a T,
	out string,
) (artifact.Local[T], error) {
	var local artifact.Local[T]

	if err := checkIsDirectory(out); err != nil {
		return local, err
	}

	m, err := mirror.Select(ctx, availableMirrors[T](), a)
	if err != nil {
		return local, err
	}

	log.Infof("downloading '%s' from mirror: %s", a.Name(), m.Name())

	remote, err := m.Remote(a)
	if err != nil {
		return local, err
	}

	c := client.NewWithRedirectDomains(m.Hosts()...)

	out = filepath.Join(out, remote.Artifact.Name())

	p, ok := ctx.Value(progressKey[T]{}).(*progress.Progress)
	if ok && p != nil {
		ctx = client.WithProgress(ctx, p)
	}

	if err := c.DownloadTo(ctx, remote.URL, out); err != nil {
		return local, err
	}

	log.Debugf("downloaded artifact: %s", out)

	local.Artifact = remote.Artifact
	local.Path = out

	return local, nil
}

/* -------------------------------------------------------------------------- */
/*                         Function: availableMirrors                         */
/* -------------------------------------------------------------------------- */

// availableMirrors returns the list of possible 'Mirror' hosts.
func availableMirrors[T artifact.Artifact]() []mirror.Mirror[T] {
	return []mirror.Mirror[T]{mirror.GitHub[T]{}, mirror.TuxFamily[T]{}}
}

/* -------------------------------------------------------------------------- */
/*                         Function: checkIsDirectory                         */
/* -------------------------------------------------------------------------- */

// checkIsDirectory is a convenience function which returns whether the provided
// path is a directory.
func checkIsDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: expected a directory: '%s'", fs.ErrInvalid, path)
	}

	return nil
}
