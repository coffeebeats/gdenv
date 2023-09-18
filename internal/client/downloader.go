package client

import (
	"io"

	"github.com/coffeebeats/gdenv/internal/progress"
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Downloader                           */
/* -------------------------------------------------------------------------- */

// An interface specifying simple methods for downloading files.
type Downloader[S any] interface {
	Download(source S, w ...io.Writer) error
	DownloadTo(source S, out string) error
	DownloadToWithProgress(source S, out string, p *progress.Progress) error
}
