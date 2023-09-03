package mirror

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/client"
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
var _ Mirror = &GitHub{} //nolint:exhaustruct

/* --------------------------- Function: NewGitHub -------------------------- */

// Creates a new GitHub 'Mirror' client with default retry mechanisms and
// redirect policies configured.
func NewGitHub() GitHub {
	client := client.Default()

	// Allow redirects to the GitHub content domain.
	client.AllowRedirectsTo(gitHubContentDomain)

	return GitHub{client}
}

/* ---------------------------- Method: Checksum ---------------------------- */

// Returns an 'Asset' to download the checksums file for the specified version
// from GitHub.
func (m *GitHub) Checksum(v godot.Version) (Asset, error) {
	var asset Asset

	if !m.Supports(v) {
		return asset, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	urlAsset, err := url.JoinPath(urlRelease, filenameChecksums)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlAsset)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	asset.name, asset.url = filenameChecksums, urlParsed

	return asset, nil
}

/* --------------------------- Method: Executable --------------------------- */

// Returns an 'Asset' to download a Godot executable for the specified version
// from GitHub.
func (m *GitHub) Executable(ex godot.Executable) (Asset, error) {
	var asset Asset

	if !m.Supports(ex.Version) {
		return asset, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, ex.Version.String())
	}

	name, err := ex.Name()
	if err != nil {
		return asset, errors.Join(ErrInvalidSpecification, err)
	}

	urlRelease, err := urlGitHubRelease(ex.Version)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	filename := name + ".zip"

	urlAsset, err := url.JoinPath(urlRelease, filename)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlAsset)
	if err != nil {
		return asset, errors.Join(ErrInvalidURL, err)
	}

	asset.name, asset.url = filename, urlParsed

	return asset, nil
}

/* ------------------------------- Method: Has ------------------------------ */

// Returns whether the mirror supports the specified version. This does *NOT*
// guarantee that the mirror has the version.
func (m *GitHub) Supports(v godot.Version) bool {
	return v.IsStable() && v.CompareNormal(versionGitHubAssetSupport) >= 0
}

/* ----------------------- Function: urlGitHubRelease ----------------------- */

// Returns a URL to the version-specific release containing release assets.
func urlGitHubRelease(v godot.Version) (string, error) {
	tag := fmt.Sprintf("%s-%s", v.Normal(), godot.LabelStable)
	return url.JoinPath(gitHubAssetsURLBase, tag)
}
