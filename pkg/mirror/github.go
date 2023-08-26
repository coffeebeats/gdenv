package mirror

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

const (
	gitHubContentDomain = "objects.githubusercontent.com"
	gitHubAssetsURLBase = "https://github.com/godotengine/godot/releases/download"
)

var versionGitHubAssetSupport = godot.MustParseVersion("v3.1.1")

/* -------------------------------------------------------------------------- */
/*                               Struct: GitHub                               */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via releases on the Godot
// GitHub repository.
type GitHub struct {
	client *resty.Client
}

/* --------------------------- Function: NewGitHub -------------------------- */

// Creates a new GitHub 'Mirror' client with default retry mechanisms and
// redirect policies configured.
func NewGitHub() GitHub {
	client := newClient()

	// Allow redirects to the GitHub content domain.
	client.SetRedirectPolicy(resty.DomainCheckRedirectPolicy(gitHubContentDomain))

	return GitHub{client}
}

/* ---------------------------- Method: Checksum ---------------------------- */

// Returns an 'Asset' to download the checksums file for the specified version
// from GitHub.
func (m *GitHub) Checksum(v godot.Version) (Asset, error) {
	var a Asset

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	urlAsset, err := url.JoinPath(urlRelease, filenameChecksums)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlAsset)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.client, a.name, a.url = m.client, filenameChecksums, urlParsed

	return a, nil
}

/* --------------------------- Method: Executable --------------------------- */

// Returns an 'Asset' to download a Godot executable for the specified version
// from GitHub.
func (m *GitHub) Executable(ex godot.Executable) (Asset, error) {
	var a Asset

	if !m.Supports(ex.Version) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, ex.Version.String())
	}

	name, err := ex.Name()
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	urlRelease, err := urlGitHubRelease(ex.Version)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	filename := name + ".zip"

	urlAsset, err := url.JoinPath(urlRelease, filename)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlAsset)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.client, a.name, a.url = m.client, filename, urlParsed

	return a, nil
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
