package mirror

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/artifacttest"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ------------------------- Test: TuxFamily.Remote ------------------------- */

func TestTuxFamilyRemote(t *testing.T) {
	tests := []struct {
		artifact artifact.Versioned
		name     string

		url *url.URL
		err error
	}{
		// Invalid inputs
		{artifact: artifacttest.MockArtifact{}, err: ErrUnsupportedArtifact},

		// Valid inputs
		{
			artifact: executable.Archive{Artifact: executable.MustParse("Godot_v4.1.1-stable_mono_linux.x86_64")},
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1.1/mono/Godot_v4.1.1-stable_mono_linux_x86_64.zip"),
		},
		{
			artifact: executable.Archive{Artifact: executable.MustParse("Godot_v4.1-stable_linux.x86_64")},
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1/Godot_v4.1-stable_linux.x86_64.zip"),
		},
		{
			artifact: executable.Archive{Artifact: executable.MustParse("Godot_v4.0-dev.20220118_win64.exe")},
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.0/pre-alpha/4.0-dev.20220118/Godot_v4.0-dev.20220118_win64.exe.zip"),
		},
		{
			artifact: source.Archive{Artifact: source.New(version.MustParse("4.1.1-stable"))},
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1.1/godot-4.1.1-stable.tar.xz"),
		},
		{
			artifact: source.Archive{Artifact: source.New(version.MustParse("4.1.0-stable"))},
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1/godot-4.1-stable.tar.xz"),
		},
		{
			artifact: mustMakeNewExecutableChecksum(t, version.MustParse("4.1.1-stable")),
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1.1/SHA512-SUMS.txt"),
		},
		{
			artifact: mustMakeNewExecutableChecksum(t, version.MustParse("4.1.0-stable")),
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1/SHA512-SUMS.txt"),
		},
		{
			artifact: mustMakeNewExecutableChecksum(t, version.MustParse("4.0-dev.20220118")),
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.0/pre-alpha/4.0-dev.20220118/SHA512-SUMS.txt"),
		},
		{
			artifact: mustMakeNewSourceChecksum(t, version.MustParse("4.1.1-stable")),
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1.1/godot-4.1.1-stable.tar.xz.sha256"),
		},
		{
			artifact: mustMakeNewSourceChecksum(t, version.MustParse("4.1.0-stable")),
			url:      mustParseURL(t, tuxFamilyAssetsURLBase+"/4.1/godot-4.1-stable.tar.xz.sha256"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.artifact.Name(), func(t *testing.T) {
			got, err := (&TuxFamily[artifact.Versioned]{}).Remote(tc.artifact)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if got := got.Artifact; got != nil && got.Name() != tc.artifact.Name() {
				t.Errorf("output: got %v, want %v", got, tc.artifact.Name())
			}
			if got := got.URL; !reflect.DeepEqual(got, tc.url) {
				t.Errorf("output: got %v, want %v", got, tc.url)
			}
		})
	}
}
