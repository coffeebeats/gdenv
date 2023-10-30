package download

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/progress"
)

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

/* -------------------------------------------------------------------------- */
/*                         Function: downloadArtifact                         */
/* -------------------------------------------------------------------------- */

// downloadArtifact downloads an artifact and reports progress to the progress
// reporter extracted from the context using the provided key.
func downloadArtifact[T artifact.Artifact](
	ctx context.Context,
	c *client.Client,
	a artifact.Remote[T],
	out string,
	progressKey any,
) error {
	p, ok := ctx.Value(progressKey).(*progress.Progress)
	if ok && p != nil {
		ctx = client.WithProgress(ctx, p)
	}

	return c.DownloadTo(ctx, a.URL, out)
}
