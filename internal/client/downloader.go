package client

import (
	"io"

	"github.com/coffeebeats/gdenv/internal/progress"
)

/* -------------------------------------------------------------------------- */
/*                          Interface: FileDownloader                         */
/* -------------------------------------------------------------------------- */

// An interface specifying simple methods for downloading files.
type FileDownloader[S any] interface {
	Download(source S, w ...io.Writer) error
	DownloadTo(source S, out string) error
	DownloadToWithProgress(source S, out string, p *progress.Progress) error
}
