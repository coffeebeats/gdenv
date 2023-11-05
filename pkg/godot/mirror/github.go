package mirror

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	gitHubContentDomain = "objects.githubusercontent.com"
	gitHubAssetsURLBase = "https://github.com/godotengine/godot-builds/releases/download"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: GitHub                               */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via releases on the Godot
// GitHub repository.
type GitHub[T artifact.Artifact] struct{}

// Validate at compile-time that 'GitHub' implements 'Mirror' interfaces.
var _ Hoster = (*GitHub[artifact.Artifact])(nil)
var _ Remoter[artifact.Artifact] = (*GitHub[artifact.Artifact])(nil)

/* ------------------------------ Impl: Hoster ------------------------------ */

// Hosts returns the host URLs at which artifacts are hosted.
func (m GitHub[T]) Hosts() []string {
	return []string{gitHubContentDomain}
}

/* ------------------------------ Impl: Remoter ----------------------------- */

// Remote returns an 'artifact.Remote' wrapper around a specified artifact. The
// remote wrapper contains the URL at which the artifact can be downloaded.
func (m GitHub[T]) Remote(a T) (artifact.Remote[T], error) {
	var remote artifact.Remote[T]

	switch any(a).(type) { // FIXME: https://github.com/golang/go/issues/45380
	case executable.Archive, executable.Checksums:
	case source.Archive, source.Checksums:
	default:
		return remote, fmt.Errorf("%w: %T", ErrUnsupportedArtifact, a)
	}

	urlRelease := urlGitHubRelease(a.Version())

	urlParsed, err := client.ParseURL(urlRelease, a.Name())
	if err != nil {
		return remote, errors.Join(ErrInvalidURL, err)
	}

	remote.Artifact, remote.URL = a, urlParsed

	return remote, nil
}

/* ------------------------------ Impl: Mirror ------------------------------ */

// Name returns the display name of the mirror.
func (m GitHub[T]) Name() string {
	return "GitHub (github.com/godotengine/godot-builds)"
}

/* ----------------------- Function: urlGitHubRelease ----------------------- */

// Returns a URL to the version-specific release containing release assets.
func urlGitHubRelease(v version.Version) string {
	// The release will be tagged as the "normal version", but a patch version
	// of '0' will be dropped.
	var normal string

	switch v.Patch() {
	case 0:
		normal = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
	default:
		normal = v.Normal()
	}

	tag := fmt.Sprintf("%s-%s", normal, version.LabelStable)

	releaseURL, err := url.JoinPath(gitHubAssetsURLBase, tag)
	if err != nil {
		panic(err) // This indicates an error in the asset URL base constant.
	}

	return releaseURL
}
