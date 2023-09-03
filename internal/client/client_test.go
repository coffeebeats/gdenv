package client

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/progress"
	"github.com/jarcoal/httpmock"
)

/* -------------------------- Test: Client.Download ------------------------- */

func TestClientDownload(t *testing.T) {
	// Given: The name of an asset to download.
	name := "asset.zip"

	// Given: A URL hosting the asset to download.
	u, err := url.Parse("https://www.example.com/" + name)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	// Given: A temporary file to write the asset to.
	f, err := os.Create(filepath.Join(t.TempDir(), name))
	defer f.Close()

	// Given: A default 'Client' instance.
	c := Default()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder("GET", u.String(),
		httpmock.NewStringResponder(200, want))

	// When: The file is downloaded.
	if err := c.Download(u, f); err != nil {
		t.Fatalf("err: got %#v, want %#v", err, nil)
	}

	// Then: The target file should have the correct contents.
	got, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if string(got) != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}
}

/* ------------------------- Test: Client.DownloadTo ------------------------ */

func TestClientDownloadTo(t *testing.T) {
	// Given: The name of an asset to download.
	name := "asset.zip"

	// Given: A URL hosting the asset to download.
	u, err := url.Parse("https://www.example.com/" + name)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	// Given: A temporary file to write the asset to.
	f := filepath.Join(t.TempDir(), name)

	// Given: A default 'Client' instance.
	c := Default()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder("GET", u.String(),
		httpmock.NewStringResponder(200, want))

	// When: The file is downloaded.
	if err := c.DownloadTo(u, f); err != nil {
		t.Fatalf("err: got %#v, want %#v", err, nil)
	}

	// Then: The target file should have the correct contents.
	got, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if string(got) != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}
}

/* ---------------------- Test: DownloadToWithProgress ---------------------- */

func TestClientDownloadToWithProgress(t *testing.T) {
	// Given: The name of an asset to download.
	name := "asset.zip"

	// Given: A URL hosting the asset to download.
	u, err := url.Parse("https://www.example.com/" + name)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	// Given: A pointer to write progress to.
	p, err := progress.New(uint64(len([]byte(name))))
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	// Given: A temporary file to write the asset to.
	f := filepath.Join(t.TempDir(), name)

	// Given: A default 'Client' instance.
	c := Default()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder("GET", u.String(),
		httpmock.NewStringResponder(200, want).SetContentLength())

	// When: The file is downloaded.
	if err := c.DownloadToWithProgress(u, f, &p); err != nil {
		t.Fatalf("err: got %#v, want %#v", err, nil)
	}

	// Then: The target file should have the correct contents.
	got, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if string(got) != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}

	// Then: The progress value should be 100%.
	percentage, err := p.Percentage()
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if want := 1.0; percentage != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}
}
