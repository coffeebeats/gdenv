package mirror

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ------------------------------ Test: Select ------------------------------ */

func TestSelect(t *testing.T) {
	tests := []struct {
		name    string
		mirrors []Mirror[executable.Archive]
		v       version.Version
		p       platform.Platform
		expects map[string]httpmock.Responder

		want Mirror[executable.Archive]
		err  error
	}{
		// Invalid inputs
		{
			name: "no mirrors results in an error",

			err: ErrMissingMirrors,
		},
		{
			name:    "no mirror supports version",
			mirrors: []Mirror[executable.Archive]{TuxFamily[executable.Archive]{}, GitHub[executable.Archive]{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot-builds/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(400, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                          httpmock.NewBytesResponder(400, nil),
			},

			err: ErrNotFound,
		},

		// Valid inputs
		{
			name:    "one valid mirror is selected",
			mirrors: []Mirror[executable.Archive]{GitHub[executable.Archive]{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot-builds/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(200, nil),
			},

			want: GitHub[executable.Archive]{},
		},
		{
			name:    "best mirror is selected",
			mirrors: []Mirror[executable.Archive]{TuxFamily[executable.Archive]{}, GitHub[executable.Archive]{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot-builds/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(200, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                          httpmock.NewBytesResponder(200, nil),
			},

			want: TuxFamily[executable.Archive]{}, // Appears first in 'mirrors'.
		},
		{
			name:    "worse mirror is selected if best isn't available",
			mirrors: []Mirror[executable.Archive]{GitHub[executable.Archive]{}, TuxFamily[executable.Archive]{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot-builds/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(400, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                          httpmock.NewBytesResponder(200, nil),
			},

			want: TuxFamily[executable.Archive]{}, // Only mirror with successful response.
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given: A new 'Client' instance without retries.
			c := client.New()

			// Given: Retry behavior is disabled for the client.
			c.RestyClient().SetRetryCount(0)

			// Given: The 'Client' instance is assigned a mock environment.
			httpmock.ActivateNonDefault(c.RestyClient().GetClient())
			defer httpmock.DeactivateAndReset()

			for urlRaw, responder := range tc.expects {
				httpmock.RegisterResponder(resty.MethodHead, urlRaw, responder)
			}

			// Given: A 'context.Context' with the stubbed client injected.
			ctx := context.WithValue(context.Background(), clientKey{}, c)

			// When: A 'Mirror' is selected from the list of options.
			got, err := Select(ctx, tc.mirrors, executable.Archive{Inner: executable.New(tc.v, tc.p)})

			// Then: The resulting error matches expectations.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			// Then: The resulting 'Mirror' matches the expected value.
			if got != tc.want {
				t.Errorf("output: got %T, want %T", got, tc.want)
			}
		})
	}
}

/* ----------------- Function: mustMakeNewExecutableChecksum ---------------- */

func mustMakeNewExecutableChecksum(t *testing.T, v version.Version) executable.Checksums {
	c, err := executable.NewChecksums(v)
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}

	return c
}

/* ------------------- Function: mustMakeNewSourceChecksum ------------------ */

func mustMakeNewSourceChecksum(t *testing.T, v version.Version) source.Checksums {
	c, err := source.NewChecksums(v)
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}

	return c
}

/* ------------------------- Function: mustParseURL ------------------------- */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
