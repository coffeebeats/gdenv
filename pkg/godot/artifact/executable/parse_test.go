package executable

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* -------------------------- Test: ParseExecutable ------------------------- */

func TestParseExecutable(t *testing.T) {
	var (
		v1     = version.MustParse("1")
		v2     = version.MustParse("2.0-beta10")
		v3     = version.MustParse("3.0.4-alpha1")
		v4     = version.MustParse("4.0.11-dev.20230101")
		v4Mono = version.MustParse("4.0-stable_mono")
	)

	tests := []struct {
		s    string
		want Executable
		err  error
	}{
		// Invalid inputs
		{s: "", want: Executable{}, err: ErrMissingName},
		{s: "Godot-v1.0-stable-x11.32", want: Executable{}, err: ErrInvalidName},
		{s: "Godot-v1.0_stable-x11.32", want: Executable{}, err: ErrInvalidName},
		{s: "Godot_invalid_x11.32", want: Executable{}, err: version.ErrInvalid},
		{s: "Godot_v1.0-stable_invalid", want: Executable{}, err: platform.ErrUnrecognizedPlatform},

		// Valid inputs
		// Linux
		{s: "Godot_v1.0-stable_x11.32", want: Executable{v1, linux32()}},
		{s: "Godot_v2.0-beta10_x11.64", want: Executable{v2, linux64()}},
		{s: "Godot_v3.0.4-alpha1_x11.32", want: Executable{v3, linux32()}},
		{s: "Godot_v4.0.11-dev.20230101_x11.64", want: Executable{v4, linux64()}},
		{s: "Godot_v4.0-stable_mono_linux_x86_64", want: Executable{v4Mono, linux64()}},

		// Darwin
		{s: "Godot_v1.0-stable_osx.fat", want: Executable{v1, macOSUniversal()}},
		{s: "Godot_v2.0-beta10_osx.64", want: Executable{v2, macOSX86_64()}},
		{s: "Godot_v3.0.4-alpha1_osx.universal", want: Executable{v3, macOSUniversal()}},
		{s: "Godot_v4.0.11-dev.20230101_macos.universal", want: Executable{v4, macOSUniversal()}},
		{s: "Godot_v4.0-stable_mono_macos.universal", want: Executable{v4Mono, macOSUniversal()}},

		// Windows
		{s: "Godot_v1.0-stable_win32", want: Executable{v1, windows32()}},
		{s: "Godot_v2.0-beta10_win64", want: Executable{v2, windows64()}},
		{s: "Godot_v3.0.4-alpha1_win32", want: Executable{v3, windows32()}},
		{s: "Godot_v4.0.11-dev.20230101_win64", want: Executable{v4, windows64()}},
		{s: "Godot_v4.0-stable_mono_win64", want: Executable{v4Mono, windows64()}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, tc.s), func(t *testing.T) {
			got, err := Parse(tc.s)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}

}
