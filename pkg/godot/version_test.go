package godot

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"unicode"
)

/* -------------------------- Test: Version.String -------------------------- */

func TestVersionString(t *testing.T) {
	tests := []struct {
		v    Version
		want string
	}{
		// Default value
		{Version{}, "v0.0-" + labelDefault},

		// Default label
		{Version{major: 1}, "v1.0-" + labelDefault},
		{Version{major: 1, minor: 1}, "v1.1-" + labelDefault},
		{Version{major: 1, minor: 1, patch: 1}, "v1.1.1-" + labelDefault},

		{Version{minor: 1}, "v0.1-" + labelDefault},
		{Version{minor: 1, patch: 1}, "v0.1.1-" + labelDefault},

		{Version{patch: 1}, "v0.0.1-" + labelDefault},

		// Specific label
		{Version{label: "label"}, "v0.0-label"},

		{Version{major: 1, label: "label"}, "v1.0-label"},
		{Version{major: 1, minor: 1, label: "label"}, "v1.1-label"},
		{Version{major: 1, minor: 1, patch: 1, label: "label"}, "v1.1.1-label"},

		{Version{minor: 1, label: "label"}, "v0.1-label"},
		{Version{minor: 1, patch: 1, label: "label"}, "v0.1.1-label"},

		{Version{patch: 1, label: "label"}, "v0.0.1-label"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.v.String()

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* --------------------------- Test: ParseVersion --------------------------- */

func TestParseVersion(t *testing.T) {
	type test struct {
		s    string
		want Version
		err  error
	}

	// Produce one test with a valid prefix and one with an invalid prefix
	withPrefixes := func(tc []test) []test {
		out := make([]test, len(tc)*2)

		for i, t := range tc {
			// Invalid
			out[i*2] = test{fmt.Sprintf("x%s", t.s), Version{}, ErrInvalidInput}

			// Valid (depending on input)
			err := t.err
			if t.s == "" || !unicode.IsDigit(rune(t.s[0])) {
				err = ErrInvalidInput
			}

			out[i*2+1] = test{fmt.Sprintf("v%s", t.s), t.want, err}
		}

		return out
	}

	// Produce tests with varying types of valid and invalid suffixes.
	withSuffixes := func(tc []test) []test {
		out := make([]test, len(tc)*2)

		for i, t := range tc {
			// Invalid
			out[i*2] = test{fmt.Sprintf("%s-", t.s), Version{}, ErrInvalidInput}

			// Valid (depending on input)
			s := "suffix"
			want, err := Version{t.want.major, t.want.minor, t.want.patch, s}, t.err
			if t.err != nil {
				want.label = ""
				err = ErrInvalidInput
			}

			out[i*2+1] = test{fmt.Sprintf("%s-%s", t.s, s), want, err}
		}

		return out
	}

	// Define tests with base version strings. These will be mutated to also
	// include prefixed, suffixed, and prefix-and-suffixed versions.
	tests := []test{
		// Invalid inputs
		{s: "", want: Version{}, err: ErrMissingInput},

		{s: "a", want: Version{}, err: ErrInvalidInput},
		{s: "0.a", want: Version{}, err: ErrInvalidInput},
		{s: "0.0.a", want: Version{}, err: ErrInvalidInput},

		{s: "0.", want: Version{}, err: ErrInvalidInput},
		{s: "0.0.", want: Version{}, err: ErrInvalidInput},
		{s: "0.0.0.", want: Version{}, err: ErrInvalidInput},

		{s: "-0", want: Version{}, err: ErrInvalidInput},
		{s: "0.-0", want: Version{}, err: ErrInvalidInput},
		{s: "0.0.-0", want: Version{}, err: ErrInvalidInput},

		{s: "00", want: Version{}, err: ErrInvalidInput},
		{s: "0.00", want: Version{}, err: ErrInvalidInput},
		{s: "0.0.00", want: Version{}, err: ErrInvalidInput},

		// Valid inputs
		{s: "1", want: Version{major: 1}, err: nil},
		{s: "1.1", want: Version{major: 1, minor: 1}, err: nil},
		{s: "1.1.1", want: Version{major: 1, minor: 1, patch: 1}, err: nil},

		{s: "0.1", want: Version{minor: 1}, err: nil},
		{s: "0.1.1", want: Version{minor: 1, patch: 1}, err: nil},

		{s: "0.0.1", want: Version{patch: 1}, err: nil},
	}

	tests = append(
		tests,
		append(
			append(
				withPrefixes(tests),
				withSuffixes(tests)...,
			),
			withSuffixes(withPrefixes(tests))...,
		)...,
	)

	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			got, err := ParseVersion(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
