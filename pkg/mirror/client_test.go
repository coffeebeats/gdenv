package mirror

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

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

	// Given: An 'Asset' representing the file to download.
	asset := Asset{name: name, url: u}

	// Given: A temporary file to write the asset to.
	f, err := os.Create(filepath.Join(t.TempDir(), asset.Name()))
	defer f.Close()

	// Given: A default 'Client' instance.
	c := Client{defaultRestyClient()}

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := asset.Name()
	httpmock.RegisterResponder("GET", asset.URL().String(),
		httpmock.NewStringResponder(200, want))

	// When: The method under test is called.
	if err := c.Download(asset, f); err != nil {
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

	// Given: An 'Asset' representing the file to download.
	asset := Asset{name: name, url: u}

	// Given: A temporary file to write the asset to.
	f := filepath.Join(t.TempDir(), asset.Name())

	// Given: A default 'Client' instance.
	c := Client{defaultRestyClient()}

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := asset.Name()
	httpmock.RegisterResponder("GET", asset.URL().String(),
		httpmock.NewStringResponder(200, want))

	// When: The method under test is called.
	if err := c.DownloadTo(asset, f); err != nil {
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

	// Given: An 'Asset' representing the file to download.
	asset := Asset{name: name, url: u}

	// Given: A pointer to write progress to.
	var progress float64

	// Given: A temporary file to write the asset to.
	f := filepath.Join(t.TempDir(), asset.Name())

	// Given: A default 'Client' instance.
	c := Client{defaultRestyClient()}

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	want := asset.Name()
	httpmock.RegisterResponder("GET", asset.URL().String(),
		httpmock.NewStringResponder(200, want).SetContentLength())

	// When: The method under test is called.
	if err := c.DownloadToWithProgress(asset, f, &progress); err != nil {
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
	if got, want := progress, 1.0; got != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}
}

/* -------------------------- Test: ProgressWriter -------------------------- */

func TestProgressWriter(t *testing.T) {
	// Given: A float containing a progress value.
	var got float64

	// Given: A 'ProgressWriter' with a size of '4', writing to 'got'.
	w, err := NewProgressWriter(4, &got)
	if err != nil {
		t.Fatalf("err: got %#v, want %#v", err, nil)
	}

	for i, b := range []byte{1, 1, 1, 1} {
		// Given: The correct initial progress value.
		if want := float64(i) / float64(4); got != want {
			t.Fatalf("output: got %#v, want %#v", got, want)
		}

		// When: A byte is written.
		n, err := w.Write([]byte{b})
		if err != nil {
			t.Fatalf("err: got %#v, want %#v", err, nil)
		}
		if n != 1 {
			t.Fatalf("output: got %#v, want %#v", n, 1)
		}

		// Then: The progress pointer updates accordingly.
		if want := float64(i+1) / float64(4); got != want {
			t.Fatalf("output: got %#v, want %#v", got, want)
		}
	}
}
