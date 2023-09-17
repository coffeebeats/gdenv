package artifact

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

/* ----------------------------- Test: NewRemote ---------------------------- */

func TestNewRemote(t *testing.T) {
	tests := []struct {
		name, url string
		want      Remote[testArtifact]
		err       error
	}{
		// Invalid inputs
		{name: "file.zip", url: "", want: Remote[testArtifact]{}, err: ErrMissingURL},
		{name: "file.zip", url: "://invalid-", want: Remote[testArtifact]{}, err: ErrInvalidURL},

		// Valid inputs
		{
			name: "file.zip",
			url:  "https://example.com",
			want: Remote[testArtifact]{Artifact: newTestArtifact("file.zip"), URL: mustParseURL(t, "https://example.com")},
			err:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%s", tc.name, tc.url), func(t *testing.T) {
			// When: A new asset is created with the specified values.
			got, err := NewRemote(newTestArtifact(tc.name), tc.url)

			// Then: The resulting error matches expectations.
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The resulting 'Asset' matches expectations.
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Struct: testArtifact -------------------------- */

type testArtifact struct {
	name string
}

func newTestArtifact(name string) testArtifact {
	return testArtifact{name: name}
}

func (testArtifact) Name() string {
	return ""
}

/* ------------------------- Function:  mustParseURL ------------------------ */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
