package platform

import (
	"errors"
	"fmt"
	"testing"
)

/* ------------------------------ Test: ParseOS ----------------------------- */

func TestParseOS(t *testing.T) {
	tests := []struct {
		s    string
		want OS
		err  error
	}{
		// Invalid inputs
		{s: "", err: ErrMissingOS},
		{s: "abc", err: ErrUnrecognizedOS},
		{s: "linux-", err: ErrUnrecognizedOS},
		{s: "mac.os", err: ErrUnrecognizedOS},
		{s: "win32", err: ErrUnrecognizedOS},

		// Valid inputs (Go-defined)
		{s: "linux", want: Linux},

		{s: "darwin", want: MacOS},
		{s: "macos", want: MacOS},
		{s: "osx", want: MacOS},

		{s: "win", want: Windows},
		{s: "windows", want: Windows},

		// Valid inputs (user-supplied)
		{s: "LINUX", want: Linux},
		{s: " LINUX\n", want: Linux},
		{s: "\tOSX ", want: MacOS},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ParseOS(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
