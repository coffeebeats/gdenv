package update

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"

	"github.com/coffeebeats/gdenv/internal/client"
)

/* -------------------------------------------------------------------------- */
/*                           Function: TestNormalize                          */
/* -------------------------------------------------------------------------- */

func TestNormalize(t *testing.T) {
	tests := []struct{ in, want string }{
		{in: "", want: ""},
		{in: "0.7.0", want: "v0.7.0"},
		{in: "v0.7.0", want: "v0.7.0"},
		{in: "  v0.7.0  ", want: "v0.7.0"},

		// Releases are always tagged with a full 'vX.Y.Z', so a partial version
		// is completed rather than left to 404 at download time.
		{in: "0.7", want: "v0.7.0"},
		{in: "1", want: "v1.0.0"},

		// Build metadata is not part of the tag.
		{in: "v0.7.0+build", want: "v0.7.0"},

		// Input which is not valid semver must stay non-empty, so that the
		// caller rejects it instead of silently updating to the latest release.
		{in: "not-a-version", want: "vnot-a-version"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := Normalize(tc.in); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: TestIsUpgrade                          */
/* -------------------------------------------------------------------------- */

func TestIsUpgrade(t *testing.T) {
	tests := []struct {
		name            string
		current, latest string
		want            bool
	}{
		{name: "newer patch", current: "v0.6.35", latest: "v0.6.36", want: true},
		{name: "newer minor", current: "v0.6.35", latest: "v0.7.0", want: true},
		{name: "identical", current: "v0.6.35", latest: "v0.6.35", want: false},

		// GitHub resolves '/releases/latest' by publish date, so a patch cut
		// against an older branch can resolve to a lower version. It must never
		// be presented as an upgrade.
		{name: "older is never an upgrade", current: "v0.7.0", latest: "v0.6.36", want: false},

		{name: "invalid current", current: "not-a-version", latest: "v0.7.0", want: false},
		{name: "invalid latest", current: "v0.6.35", latest: "not-a-version", want: false},
		{name: "both empty", current: "", latest: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUpgrade(tc.current, tc.latest); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestLatestVersion                        */
/* -------------------------------------------------------------------------- */

func TestLatestVersion(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		location string
		want     string
		err      error
	}{
		{
			name:     "resolves the tag from the redirect",
			status:   http.StatusFound,
			location: "https://github.com/coffeebeats/gdenv/releases/tag/v0.7.0",
			want:     "v0.7.0",
		},
		{
			name:     "non-semver tag is rejected",
			status:   http.StatusFound,
			location: "https://github.com/coffeebeats/gdenv/releases/tag/nightly",
			err:      ErrInvalidVersion,
		},
		{
			name:   "missing location is reported",
			status: http.StatusFound,
			err:    client.ErrMissingLocation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := client.NewWithoutRedirects()

			httpmock.ActivateNonDefault(c.RestyClient().GetClient())
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(resty.MethodGet,
				repositoryURL+"/"+pathReleaseLatest,
				func(_ *http.Request) (*http.Response, error) {
					res := httpmock.NewStringResponse(tc.status, "")
					if tc.location != "" {
						res.Header.Set("Location", tc.location)
					}

					return res, nil
				},
			)

			got, err := latestVersion(context.Background(), c)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			if tc.err != nil {
				return
			}

			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                        Function: TestLocateBinaries                        */
/* -------------------------------------------------------------------------- */

func TestLocateBinaries(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		want    []string
		wantErr error
	}{
		{
			name:  "finds binaries at the archive root",
			files: []string{"gdenv", "godot"},
			want:  []string{"gdenv", "godot"},
		},
		{
			// The shim is bundled via a 'files:' glob in '.goreleaser.yaml', so
			// its directory prefix must not matter.
			name:  "finds binaries nested under a directory",
			files: []string{"gdenv", "gdenv-shim_linux_amd64/godot"},
			want:  []string{"gdenv", "godot"},
		},
		{
			name:    "reports a missing binary",
			files:   []string{"gdenv"},
			want:    []string{"gdenv", "godot"},
			wantErr: ErrMissingBinary,
		},
		{
			name:    "reports a duplicate binary",
			files:   []string{"a/godot", "b/godot", "gdenv"},
			want:    []string{"gdenv", "godot"},
			wantErr: ErrDuplicateBinary,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			for _, f := range tc.files {
				path := filepath.Join(dir, filepath.FromSlash(f))

				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}

				if err := os.WriteFile(path, []byte(f), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			found, err := locateBinaries(dir, tc.want)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err: got %v, want %v", err, tc.wantErr)
			}

			if tc.wantErr != nil {
				return
			}

			for _, name := range tc.want {
				if _, ok := found[name]; !ok {
					t.Fatalf("did not locate %q in %v", name, found)
				}
			}
		})
	}
}
