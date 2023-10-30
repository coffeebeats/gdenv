package version

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
	"unicode"
)

/* ------------------------------- Test: Parse ------------------------------ */

func TestParse(t *testing.T) {
	type test struct {
		s    string
		want Version
		err  error
	}

	// Produce one test with a valid prefix and one with an invalid prefix
	withPrefixes := func(tc []test) []test {
		out := make([]test, len(tc)*3)

		for i, t := range tc {
			// Invalid
			out[i*3] = test{fmt.Sprintf("x%s", t.s), Version{}, ErrInvalid}

			// Valid (depending on input)
			err := t.err
			if t.s == "" || !unicode.IsDigit(rune(t.s[0])) {
				err = ErrInvalid
			}

			out[i*3+1] = test{fmt.Sprintf("v%s", t.s), t.want, err}

			// Valid (depending on input)
			err = t.err
			if t.s == "" || !unicode.IsDigit(rune(t.s[0])) {
				err = ErrInvalid
			}

			out[i*3+2] = test{fmt.Sprintf("\t \nV%s", t.s), t.want, err}
		}

		return out
	}

	// Produce tests with varying types of valid and invalid suffixes.
	withSuffixes := func(tc []test) []test {
		out := make([]test, len(tc)*3)

		for i, t := range tc {
			// Invalid
			out[i*3] = test{fmt.Sprintf("%s-", t.s), Version{}, ErrInvalid}

			// Valid (depending on input)
			s := "suffix"
			want, err := Version{t.want.major, t.want.minor, t.want.patch, s}, t.err
			if t.err != nil {
				want.label = ""
				err = ErrInvalid
			}

			out[i*3+1] = test{fmt.Sprintf("%s-%s", t.s, s), want, err}

			// Valid (depending on input)
			sNormalized, s := s, "SUFFIX\t\n "

			want, err = Version{t.want.major, t.want.minor, t.want.patch, sNormalized}, t.err
			if t.err != nil {
				want.label = ""
				err = ErrInvalid
			}

			out[i*3+2] = test{fmt.Sprintf("%s-%s", t.s, s), want, err}
		}

		return out
	}

	// Define tests with base version strings. These will be mutated to also
	// include prefixed, suffixed, and prefix-and-suffixed versions.
	tests := []test{
		// Invalid inputs
		{s: "", want: Version{}, err: ErrMissing},

		{s: "a", want: Version{}, err: ErrInvalid},
		{s: "0.a", want: Version{}, err: ErrInvalid},
		{s: "0.0.a", want: Version{}, err: ErrInvalid},

		{s: "0.", want: Version{}, err: ErrInvalid},
		{s: "0.0.", want: Version{}, err: ErrInvalid},
		{s: "0.0.0.", want: Version{}, err: ErrInvalid},

		{s: "-0", want: Version{}, err: ErrInvalid},
		{s: "0.-0", want: Version{}, err: ErrInvalid},
		{s: "0.0.-0", want: Version{}, err: ErrInvalid},

		{s: "00", want: Version{}, err: ErrInvalid},
		{s: "0.00", want: Version{}, err: ErrInvalid},
		{s: "0.0.00", want: Version{}, err: ErrInvalid},

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
			got, err := Parse(tc.s)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------- Function: TestParseInvalidNumber -------------------- */

func TestParseInvalidNumber(t *testing.T) {
	tests := []struct {
		input string

		want Version
		err  error
	}{
		{input: "v1." + strconv.FormatUint(math.MaxUint64, 10), err: ErrInvalidNumber},
		{input: "v1." + strconv.FormatInt(math.MinInt64, 10), err: ErrInvalid},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}

}
