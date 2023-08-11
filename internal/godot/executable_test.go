package godot

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

/* -------------------------- Test: ExecutableName -------------------------- */

func TestExecutableName(t *testing.T) {
	tests := []struct {
		v      Version
		target string
		want   string
		err    error
	}{
		{Version{"0", "0", "0", "stable"}, "t", "Godot_v0.0.0-stable_t", nil},
		{Version{"0", "0", "", "stable"}, "t", "Godot_v0.0-stable_t", nil},
		{Version{"0", "", "", "stable"}, "t", "Godot_v0-stable_t", nil},

		{Version{"0", "0", "0", ""}, "t", "Godot_v0.0.0-stable_t", nil},
		{Version{"0", "0", "", ""}, "t", "Godot_v0.0-stable_t", nil},
		{Version{"0", "", "", ""}, "t", "Godot_v0-stable_t", nil},

		{Version{}, "t", "", ErrInvalidVersion},
		{Version{Minor: "1"}, "t", "", ErrInvalidVersion},
		{Version{Patch: "1"}, "t", "", ErrInvalidVersion},
		{Version{Suffix: "1"}, "t", "", ErrInvalidVersion},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if err := os.Setenv(envVarPlatform, tc.target); err != nil {
				t.Fatalf("Failed to update environment: %v", err)
			}

			got, err := ExecutableName(tc.v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Fatalf("output: got %v, want %v", got, tc.want)
			}
		})
	}
}
