package github

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/pkg/godot"
)

const (
	gitHubContentDomain = "objects.githubusercontent.com"
	gitHubAssetsURLBase = "https://github.com/godotengine/godot/releases/download"
)

var versionGitHubAssetSupport = godot.MustParseVersion("v3.1.1") //nolint:gochecknoglobals

/* -------------------------------------------------------------------------- */
/*                               Struct: GitHub                               */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via releases on the Godot
// GitHub repository.
type GitHub struct {
	client *client.Client
}

// Validate at compile-time that 'GitHub' implements 'Mirror'.
var _ mirror.Mirror = &GitHub{} //nolint:exhaustruct

/* ------------------------------ Function: New ----------------------------- */

// Creates a new GitHub 'Mirror' client with default retry mechanisms and
// redirect policies configured.
func New() GitHub {
	c := client.NewWithRedirectDomains(gitHubContentDomain)
	return GitHub{&c}
}

/* ----------------------------- Impl: Checksum ----------------------------- */

// Returns an 'Asset' to download the checksums file for the specified version
// from GitHub.
func (m GitHub) Checksum(v godot.Version) (mirror.Asset, error) {
	if !m.Supports(v) {
		return mirror.Asset{}, fmt.Errorf("%w: '%s'", mirror.ErrInvalidSpecification, v.String())
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	urlRaw, err := url.JoinPath(urlRelease, mirror.FilenameChecksums)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	return mirror.NewAsset(mirror.FilenameChecksums, urlRaw)
}

/* ---------------------------- Impl: Executable ---------------------------- */

// Returns an 'Asset' to download a Godot executable for the specified version
// from GitHub.
func (m GitHub) Executable(ex godot.Executable) (mirror.Asset, error) {
	if !m.Supports(ex.Version) {
		return mirror.Asset{}, fmt.Errorf("%w: '%s'", mirror.ErrInvalidSpecification, ex.Version.String())
	}

	name, err := ex.Name()
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidSpecification, err)
	}

	urlRelease, err := urlGitHubRelease(ex.Version)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	filename := name + ".zip"

	urlRaw, err := url.JoinPath(urlRelease, filename)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	return mirror.NewAsset(filename, urlRaw)
}

/* -------------------------------- Impl: Has ------------------------------- */

// Issues a request to see if the mirror host has the specific version.
func (m GitHub) Has(v godot.Version) bool {
	if !m.Supports(v) {
		return false
	}

	// Rather than maintaining a separate source of truth, issue a HEAD request
	// to test whether the version exists.
	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return false
	}

	exists, err := m.client.Exists(urlRelease)
	if err != nil {
		return false
	}

	return exists
}

/* ----------------------------- Impl: Supports ----------------------------- */

// Checks whether the version is broadly supported by the mirror. No network
// request is issued, but this does not guarantee the host has the version.
// To check whether the host has the version definitively via the network,
// use the 'Has' method.
func (m GitHub) Supports(v godot.Version) bool {
	// GitHub only contains stable releases, starting with 'versionGitHubAssetSupport'.
	return v.IsStable() && v.CompareNormal(versionGitHubAssetSupport) >= 0
}

/* ----------------------- Function: urlGitHubRelease ----------------------- */

// Returns a URL to the version-specific release containing release assets.
func urlGitHubRelease(v godot.Version) (string, error) {
	// The release will be tagged as the "normal version", but a patch version
	// of '0' will be dropped.
	var normal string

	switch v.Patch() {
	case 0:
		normal = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
	default:
		normal = v.Normal()
	}

	tag := fmt.Sprintf("%s-%s", normal, godot.LabelStable)

	return url.JoinPath(gitHubAssetsURLBase, tag)
}
