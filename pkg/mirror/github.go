package mirror

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

const (
	gitHubContentDomain        = "githubusercontent.com"
	gitHubReleaseAssetsURLBase = "https://github.com/godotengine/godot/releases/download"
)

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

	urlRaw, err := url.JoinPath(gitHubReleaseAssetsURLBase, v.String(), filenameChecksums)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlRaw)
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

	filename := name + ".zip"

	urlRaw, err := url.JoinPath(gitHubReleaseAssetsURLBase, ex.Version.String(), filename)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	urlParsed, err := url.Parse(urlRaw)
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
	return v.IsStable()
}
