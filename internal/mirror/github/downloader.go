package github

import (
	"io"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/internal/progress"
)

// Validate at compile-time that 'GitHub' implements 'FileDownloader'.
var _ client.FileDownloader[mirror.Asset] = &GitHub{} //nolint:exhaustruct

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (m *GitHub) Download(a mirror.Asset, w ...io.Writer) error {
	return m.client.Download(a.URL(), w...)
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (m *GitHub) DownloadTo(a mirror.Asset, out string) error {
	return m.client.DownloadTo(a.URL(), out)
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (m *GitHub) DownloadToWithProgress(a mirror.Asset, out string, p *progress.Progress) error {
	return m.client.DownloadToWithProgress(a.URL(), out, p)
}
