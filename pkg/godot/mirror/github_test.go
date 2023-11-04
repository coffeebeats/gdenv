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

/* --------------------------- Test: GitHub.Remote -------------------------- */

func TestGitHubRemote(t *testing.T) {
	tests := []struct {
		artifact artifact.Artifact

		url *url.URL
		err error
	}{
		// Invalid inputs
		{artifact: artifacttest.MockArtifact{}, err: ErrUnsupportedArtifact},

		// Valid inputs
		{
			artifact: executable.Archive{Inner: executable.MustParse("Godot_v4.1.1-stable_linux.x86_64")},
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1.1-stable/Godot_v4.1.1-stable_linux.x86_64.zip"),
		},
		{
			artifact: executable.Archive{Inner: executable.MustParse("Godot_v4.1-stable_linux.x86_64")},
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1-stable/Godot_v4.1-stable_linux.x86_64.zip"),
		},
		{
			artifact: source.Archive{Inner: source.New(version.MustParse("4.1.1-stable"))},
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1.1-stable/godot-4.1.1-stable.tar.xz"),
		},
		{
			artifact: source.Archive{Inner: source.New(version.MustParse("4.1.0-stable"))},
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1-stable/godot-4.1-stable.tar.xz"),
		},
		{
			artifact: mustMakeNewExecutableChecksum(t, version.MustParse("4.1.1-stable")),
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1.1-stable/SHA512-SUMS.txt"),
		},
		{
			artifact: mustMakeNewExecutableChecksum(t, version.MustParse("4.1.0-stable")),
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1-stable/SHA512-SUMS.txt"),
		},
		{
			artifact: mustMakeNewSourceChecksum(t, version.MustParse("4.1.1-stable")),
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1.1-stable/godot-4.1.1-stable.tar.xz.sha256"),
		},
		{
			artifact: mustMakeNewSourceChecksum(t, version.MustParse("4.1.0-stable")),
			url:      mustParseURL(t, gitHubAssetsURLBase+"/4.1-stable/godot-4.1-stable.tar.xz.sha256"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.artifact.Name(), func(t *testing.T) {
			got, err := (&GitHub[artifact.Artifact]{}).Remote(tc.artifact)

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
