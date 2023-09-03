package client

import (
	"io"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/progress"
)

/* -------------------------------------------------------------------------- */
/*                          Interface: FileDownloader                         */
/* -------------------------------------------------------------------------- */

// An interface specifying simple methods for downloading files.
type FileDownloader interface {
	Download(u *url.URL, w ...io.Writer) error
	DownloadTo(u *url.URL, out string) error
	DownloadToWithProgress(u *url.URL, out string, p *progress.Progress) error
}
