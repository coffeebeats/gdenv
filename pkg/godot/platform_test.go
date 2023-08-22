package godot

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
		{s: "", err: ErrMissingTargetInput},
		{s: "abc", err: ErrUnrecognizedTargetInput},
		{s: "linux-", err: ErrUnrecognizedTargetInput},
		{s: "mac.os", err: ErrUnrecognizedTargetInput},
		{s: "win32", err: ErrUnrecognizedTargetInput},

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

/* ----------------------------- Test: ParseArch ---------------------------- */

func TestParseArch(t *testing.T) {
	tests := []struct {
		s    string
		want Arch
		err  error
	}{
		// Invalid inputs
		{s: "", err: ErrMissingTargetInput},
		{s: "abc", err: ErrUnrecognizedTargetInput},

		// Valid inputs (Go-defined)
		{s: "amd64", want: Amd64},
		{s: "x86_64", want: Amd64},
		{s: "x86-64", want: Amd64},

		{s: "arm64", want: Arm64},
		{s: "arm64be", want: Arm64},

		{s: "386", want: I386},
		{s: "i386", want: I386},
		{s: "x86", want: I386},

		{s: "fat", want: Universal},
		{s: "universal", want: Universal},

		// Valid inputs (user-supplied)
		{s: "AMD64", want: Amd64},
		{s: " X86_64\n", want: Amd64},
		{s: "\tuniversal ", want: Universal},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ParseArch(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
