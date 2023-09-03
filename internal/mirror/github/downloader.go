package github

import (
	"io"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/progress"
)

// Validate at compile-time that 'GitHub' implements 'FileDownloader'.
var _ client.FileDownloader = &GitHub{} //nolint:exhaustruct

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (m *GitHub) Download(u *url.URL, w ...io.Writer) error {
	return m.client.Download(u, w...)
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (m *GitHub) DownloadTo(u *url.URL, out string) error {
	return m.client.DownloadTo(u, out)
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (m *GitHub) DownloadToWithProgress(u *url.URL, out string, p *progress.Progress) error {
	return m.client.DownloadToWithProgress(u, out, p)
}
