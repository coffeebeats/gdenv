package godot

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"unicode"
)

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
		{input: "0", want: Version{"0", "0", "0", "stable"}, err: nil},
		{input: "0.0", want: Version{"0", "0", "0", "stable"}, err: nil},
		{input: "0.0.0", want: Version{"0", "0", "0", "stable"}, err: nil},
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

			if err != tc.err && errors.Unwrap(err) != tc.err {
				fmt.Println(errors.Unwrap(err))
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
