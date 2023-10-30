package mirror

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

/* ------------------------------ Test: Select ------------------------------ */

func TestSelect(t *testing.T) {
	tests := []struct {
		name    string
		mirrors []Mirror
		v       version.Version
		p       platform.Platform
		expects map[string]httpmock.Responder

		want Mirror
		err  error
	}{
		// Invalid inputs
		{
			name: "no mirrors results in an error",

			err: ErrMissingMirrors,
		},
		{
			name:    "no mirror supports version",
			mirrors: []Mirror{TuxFamily{}, GitHub{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(400, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                   httpmock.NewBytesResponder(400, nil),
			},

			err: ErrNotFound,
		},

		// Valid inputs
		{
			name:    "one valid mirror is selected",
			mirrors: []Mirror{GitHub{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(200, nil),
			},

			want: GitHub{},
		},
		{
			name:    "best mirror is selected",
			mirrors: []Mirror{TuxFamily{}, GitHub{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(200, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                   httpmock.NewBytesResponder(200, nil),
			},

			want: TuxFamily{}, // Appears first in 'mirrors'.
		},
		{
			name:    "worse mirror is selected if best isn't available",
			mirrors: []Mirror{GitHub{}, TuxFamily{}},
			v:       version.Godot4(),
			p:       platform.MustParse("win64"),
			expects: map[string]httpmock.Responder{
				"https://github.com/godotengine/godot/releases/download/4.0-stable/Godot_v4.0-stable_win64.exe.zip": httpmock.NewBytesResponder(400, nil),
				"https://downloads.tuxfamily.org/godotengine/4.0/Godot_v4.0-stable_win64.exe.zip":                   httpmock.NewBytesResponder(200, nil),
			},

			want: TuxFamily{}, // Only mirror with successful response.
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
			got, err := Select(ctx, tc.v, tc.p, tc.mirrors)

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

/* ----------------------- Function: mustMakeNewSource ---------------------- */

func mustMakeNewSource(t *testing.T, v version.Version) checksum.Source {
	t.Helper()

	s, err := checksum.NewSource(v)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return s
}

/* ------------------------- Function: mustParseURL ------------------------- */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
