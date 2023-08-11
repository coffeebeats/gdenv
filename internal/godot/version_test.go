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
		{Version{}, ""},

		{Version{"", "1", "", ""}, ""},
		{Version{"", "", "1", ""}, ""},
		{Version{"", "", "", "s"}, ""},

		{Version{"1", "", "", ""}, "v1"},
		{Version{"1", "1", "", ""}, "v1.1"},
		{Version{"1", "1", "1", ""}, "v1.1.1"},

		{Version{"1", "", "", "s"}, "v1-s"},
		{Version{"1", "1", "", "s"}, "v1.1-s"},
		{Version{"1", "1", "1", "s"}, "v1.1.1-s"},
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

/* ------------------------- Test: Version.Canonical ------------------------ */

func TestVersionCanonical(t *testing.T) {
	tests := []struct {
		v    Version
		want string
	}{
		{Version{}, "v0.0.0-stable"},

		{Version{"1", "", "", ""}, "v1.0.0-stable"},
		{Version{"1", "1", "", ""}, "v1.1.0-stable"},
		{Version{"1", "1", "1", ""}, "v1.1.1-stable"},

		{Version{"1", "", "", "s"}, "v1.0.0-s"},
		{Version{"1", "1", "", "s"}, "v1.1.0-s"},
		{Version{"1", "1", "1", "s"}, "v1.1.1-s"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.v.Canonical()

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Test: Version.IsValid ------------------------- */

func TestVersionIsValid(t *testing.T) {
	tests := []struct {
		v    Version
		want bool
	}{
		{Version{}, false},

		{Version{Minor: "1"}, false},
		{Version{Patch: "1"}, false},
		{Version{Suffix: "s"}, false},
		{Version{Minor: "1", Patch: "1"}, false},
		{Version{Minor: "1", Suffix: "s"}, false},
		{Version{Patch: "1", Suffix: "s"}, false},
		{Version{Minor: "1", Patch: "1", Suffix: "s"}, false},

		{Version{"1", "", "", ""}, true},
		{Version{"1", "1", "", ""}, true},
		{Version{"1", "1", "1", ""}, true},

		{Version{"1", "", "", "s"}, true},
		{Version{"1", "1", "", "s"}, true},
		{Version{"1", "1", "1", "s"}, true},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.v.IsValid()

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* --------------------------- Test: ParseVersion --------------------------- */

func TestParseVersion(t *testing.T) {
	type test struct {
		input string
		want  Version
		err   error
	}

	// Produce one test with a valid prefix and one with an invalid prefix
	withPrefixes := func(tc []test) []test {
		out := make([]test, len(tc)*2)

		for i, t := range tc {
			// Invalid
			out[i*2] = test{fmt.Sprintf("x%s", t.input), Version{}, ErrInvalidInput}

			// Valid (depending on input)
			err := t.err
			if t.input == "" || !unicode.IsDigit(rune(t.input[0])) {
				err = ErrInvalidInput
			}

			out[i*2+1] = test{fmt.Sprintf("v%s", t.input), t.want, err}
		}

		return out
	}

	// Produce tests with varying types of valid and invalid suffixes.
	withSuffixes := func(tc []test) []test {
		out := make([]test, len(tc)*2)

		for i, t := range tc {
			// Invalid
			out[i*2] = test{fmt.Sprintf("%s-", t.input), Version{}, ErrInvalidInput}

			// Valid (depending on input)
			s := "suffix"
			want, err := Version{t.want.Major, t.want.Minor, t.want.Patch, s}, t.err
			if t.err != nil {
				want.Suffix = ""
				err = ErrInvalidInput
			}

			out[i*2+1] = test{fmt.Sprintf("%s-%s", t.input, s), want, err}
		}

		return out
	}

	// Define tests with base version strings. These will be mutated to also
	// include prefixed, suffixed, and prefix-and-suffixed versions.
	tests := []test{
		// Invalid inputs
		{input: "", want: Version{}, err: ErrNoInput},

		{input: "a", want: Version{}, err: ErrInvalidInput},
		{input: "0.a", want: Version{}, err: ErrInvalidInput},
		{input: "0.0.a", want: Version{}, err: ErrInvalidInput},

		{input: "0.", want: Version{}, err: ErrInvalidInput},
		{input: "0.0.", want: Version{}, err: ErrInvalidInput},
		{input: "0.0.0.", want: Version{}, err: ErrInvalidInput},

		{input: "-0", want: Version{}, err: ErrInvalidInput},
		{input: "0.-0", want: Version{}, err: ErrInvalidInput},
		{input: "0.0.-0", want: Version{}, err: ErrInvalidInput},

		{input: "00", want: Version{}, err: ErrInvalidInput},
		{input: "0.00", want: Version{}, err: ErrInvalidInput},
		{input: "0.0.00", want: Version{}, err: ErrInvalidInput},

		// Valid inputs
		{input: "0", want: Version{Major: "0"}, err: nil},
		{input: "0.0", want: Version{"0", "0", "", ""}, err: nil},
		{input: "0.0.0", want: Version{"0", "0", "0", ""}, err: nil},
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
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseVersion(tc.input)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
