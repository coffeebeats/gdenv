package mirror

import (
	"net/url"
	"testing"
)

/* ------------------------- Function: mustParseURL ------------------------- */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
