package artifact

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

/* -------------------------- Test: RemoteParseURL -------------------------- */

func TestRemoteParseURL(t *testing.T) {
	tests := []struct {
		url  string
		want *url.URL
		err  error
	}{
		// Invalid inputs
		{url: "", want: nil, err: ErrMissingURL},
		{url: "://invalid-", want: nil, err: ErrInvalidURL},

		// Valid inputs
		{
			url:  "https://example.com",
			want: mustParseURL(t, "https://example.com"),
			err:  nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-url='%s'", i, tc.url), func(t *testing.T) {
			r := Remote[EmptyArtifact]{Artifact: EmptyArtifact{}, URL: tc.url}

			// When: A new asset is created with the specified values.
			got, err := r.ParseURL()

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

/* -------------------------  ------------------------ */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}

/* -------------------------------------------------------------------------- */
/*                            Struct: EmptyArtifact                           */
/* -------------------------------------------------------------------------- */

// An empty test 'Artifact'.
type EmptyArtifact struct {
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a EmptyArtifact) Name() string {
	return "artifact"
}
