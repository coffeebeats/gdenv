package mirror

import (
	"net/url"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

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
