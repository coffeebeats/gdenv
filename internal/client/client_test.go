package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/progress"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

/* ----------------------------- Test: ParseURL ----------------------------- */

func TestParseURL(t *testing.T) {
	tests := []struct {
		urlParts []string
		want     *url.URL
		err      error
	}{
		// Invalid inputs
		{urlParts: []string{""}, want: nil, err: ErrMissingURL},
		{urlParts: []string{"://invalid-"}, want: nil, err: ErrInvalidURL},

		// Valid inputs
		{
			urlParts: []string{"https://example.com"},
			want:     mustParseURL(t, "https://example.com/"),
			err:      nil,
		},
		{
			urlParts: []string{"https://example.com", "abc"},
			want:     mustParseURL(t, "https://example.com/abc"),
			err:      nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-url='%s'", i, strings.Join(tc.urlParts, "/")), func(t *testing.T) {
			// When: A new asset is created with the specified values.
			got, err := ParseURL(tc.urlParts[0], tc.urlParts[1:]...)

			// Then: The resulting error matches expectations.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The resulting 'Asset' matches expectations.
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

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

	// Given: A pointer to write progress to.
	p := progress.Progress{}

	// Given: A 'context.Context' with the specified progress reporter.
	ctx := WithProgress(context.Background(), &p)

	// Given: Mocked contents of the asset.
	httpmock.ActivateNonDefault(c.restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	want := name
	httpmock.RegisterResponder(resty.MethodGet, u.String(),
		httpmock.NewStringResponder(200, want).SetContentLength())

	// When: The file is downloaded.
	if err := c.Download(ctx, u, f); err != nil {
		t.Errorf("err: got %#v, want %#v", err, nil)
	}

	// Then: The target file should have the correct contents.
	got, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if string(got) != want {
		t.Errorf("output: got %#v, want %#v", got, want)
	}

	// Then: The progress value should be 100%.
	if got, want := p.Percentage(), 1.0; got != want {
		t.Errorf("output: got %#v, want %#v", got, want)
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

	// Given: A pointer to write progress to.
	p := progress.Progress{}

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

	// Given: A 'context.Context' with the specified progress reporter.
	ctx := WithProgress(context.Background(), &p)

	// When: The file is downloaded.
	if err := c.DownloadTo(ctx, u, f); err != nil {
		t.Errorf("err: got %#v, want %#v", err, nil)
	}

	// Then: The target file should have the correct contents.
	got, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}
	if string(got) != want {
		t.Errorf("output: got %#v, want %#v", got, want)
	}

	// Then: The progress value should be 100%.
	if got, want := p.Percentage(), 1.0; got != want {
		t.Errorf("output: got %#v, want %#v", got, want)
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
		{url: "", err: ErrMissingURL},
		{url: "https://www.example.com/", res: http.StatusMovedPermanently, err: ErrUnexpectedRedirect},

		// Valid inputs
		{url: "https://www.example.com/", res: http.StatusOK, want: true},
		{url: "https://www.example.com/", res: http.StatusNotFound, want: false},
		{url: "https://www.example.com/", res: http.StatusBadGateway, want: false},
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
			got, err := c.Exists(context.Background(), tc.url)

			// Then: The returned error matches expectations.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The returned existence value matches expectations.
			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------- Function: mustParseURL ------------------------- */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
