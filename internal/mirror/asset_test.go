package mirror

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

/* ----------------------------- Test: NewAsset ----------------------------- */

func TestNewAsset(t *testing.T) {
	tests := []struct {
		name, url string
		want      Asset
		err       error
	}{
		// Invalid inputs
		{name: "", url: "https://example.com", want: Asset{}, err: ErrMissingName},
		{name: "file.zip", url: "", want: Asset{}, err: ErrMissingURL},
		{name: "file.zip", url: "://invalid-", want: Asset{}, err: ErrInvalidURL},

		// Valid inputs
		{
			name: "file.zip",
			url:  "https://example.com",
			want: Asset{name: "file.zip", url: mustParseURL(t, "https://example.com")},
			err:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%s", tc.name, tc.url), func(t *testing.T) {
			// When: A new asset is created with the specified values.
			got, err := NewAsset(tc.name, tc.url)

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

/* ------------------------- Function:  mustParseURL ------------------------ */

func mustParseURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatalf("test setup: %#v", err)
	}

	return u
}
