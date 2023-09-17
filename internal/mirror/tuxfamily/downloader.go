package tuxfamily

import (
	"io"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/progress"
)

// Validate at compile-time that 'TuxFamily' implements 'FileDownloader'.
var _ client.FileDownloader[artifact.Remote[artifact.Artifact]] = &TuxFamily{} //nolint:exhaustruct

/* ----------------------------- Impl: Download ----------------------------- */

// Downloads the provided asset, copying the response to all of the provided
// 'io.Writer' writers.
func (m TuxFamily) Download(a artifact.Remote[artifact.Artifact], w ...io.Writer) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.Download(urlParsed, w...)
}

/* ---------------------------- Impl: DownloadTo ---------------------------- */

// Downloads the provided asset to a specified file 'out'.
func (m TuxFamily) DownloadTo(a artifact.Remote[artifact.Artifact], out string) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.DownloadTo(urlParsed, out)
}

/* ---------------------- Impl: DownloadToWithProgress ---------------------- */

// Downloads the response of a request to the specified filepath, reporting the
// download progress to the provided progress pointer 'p'.
func (m TuxFamily) DownloadToWithProgress(a artifact.Remote[artifact.Artifact], out string, p *progress.Progress) error {
	urlParsed, err := a.ParseURL()
	if err != nil {
		return err
	}

	return m.client.DownloadToWithProgress(urlParsed, out, p)
}
