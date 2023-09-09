package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/coffeebeats/gdenv/internal/progress"
	"github.com/go-resty/resty/v2"
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
	c := New()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder(resty.MethodGet, u.String(),
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
	c := New()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder(resty.MethodGet, u.String(),
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
	c := New()

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder(resty.MethodGet, u.String(),
		httpmock.NewStringResponder(200, want).SetContentLength())

	// When: The file is downloaded.
	if err := c.DownloadToWithProgress(u, f, p); err != nil {
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
	if got, want := p.Percentage(), 1.0; got != want {
		t.Fatalf("output: got %#v, want %#v", got, want)
	}
}

/* --------------------------- Test: Client.Exists -------------------------- */

func TestClientExists(t *testing.T) {
	tests := []struct {
		url  string
		res  int
		want bool
		err  error
	}{
		// Invalid inputs
		{url: "", err: ErrInvalidURL},
		{url: "https://www.example.com", res: http.StatusMovedPermanently, err: ErrUnexpectedRedirect},

		// Valid inputs
		{url: "https://www.example.com", res: http.StatusOK, want: true},
		{url: "https://www.example.com", res: http.StatusNotFound, want: false},
		{url: "https://www.example.com", res: http.StatusBadGateway, want: false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("url=%s,res=%s", tc.url, strconv.Itoa(tc.res)), func(t *testing.T) {
			// Given: A default 'Client' instance.
			c := New()
			c.restyClient.SetRetryCount(0) // Disable retries to speed up tests.

			// Given: A mocked response.
			httpmock.ActivateNonDefault(c.restyClient.GetClient())
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(resty.MethodHead, tc.url,
				httpmock.NewStringResponder(tc.res, ""))

			// When: The URL is checked for existence.
			got, err := c.Exists(tc.url)

			// Then: The returned error matches expectations.
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The returned existence value matches expectations.
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
