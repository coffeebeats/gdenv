package source

import (
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ---------------------------- Test: Source.Name --------------------------- */

func TestSourceName(t *testing.T) {
	tests := []struct {
		version version.Version
		want    string
	}{
		// Valid inputs
		{version: version.Version{}, want: "godot-0.0-stable"},
		{version: version.Godot3(), want: "godot-3.0-stable"},
		{version: version.MustParse("v4.1.1-dev1"), want: "godot-4.1.1-dev1"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-'%s'", i, tc.version), func(t *testing.T) {
			if got := (Source{tc.version}).Name(); got != tc.want {
				t.Errorf("output: got %v, want %v", got, tc.want)
			}
		})
	}
}
