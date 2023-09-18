package github

import (
	"io"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/progress"
)

// Validate at compile-time that 'GitHub' implements 'Downloader'.
var _ client.Downloader[artifact.Hosted] = &GitHub{} //nolint:exhaustruct

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (m GitHub) Download(a artifact.Hosted, w ...io.Writer) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.Download(urlParsed, w...)
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (m GitHub) DownloadTo(a artifact.Hosted, out string) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.DownloadTo(urlParsed, out)
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (m GitHub) DownloadToWithProgress(a artifact.Hosted, out string, p *progress.Progress) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.DownloadToWithProgress(urlParsed, out, p)
}
