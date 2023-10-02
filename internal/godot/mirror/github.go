package mirror

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	gitHubContentDomain = "objects.githubusercontent.com"
	gitHubAssetsURLBase = "https://github.com/godotengine/godot/releases/download"
)

var versionGitHubAssetSupport = version.MustParse("v3.1.1") //nolint:gochecknoglobals

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

/* ------------------------------ Function: New ----------------------------- */

// Creates a new GitHub 'Mirror' client with default retry mechanisms and
// redirect policies configured.
func NewGitHub() GitHub {
	c := client.NewWithRedirectDomains(gitHubContentDomain)
	return GitHub{&c}
}

/* ------------------------------ Impl: Mirror ------------------------------ */

// Issues a request to see if the mirror host has the specific version.
func (m GitHub) CheckIfSupports(v version.Version) bool {
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

func (m GitHub) ExecutableArchive(ex executable.Executable) (artifact.Remote[executable.Archive], error) {
	var a artifact.Remote[executable.Archive]

	if !m.Supports(ex.Version()) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, ex.Version())
	}

	urlRelease, err := urlGitHubRelease(ex.Version())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	executableArchive := executable.Archive{Artifact: ex}

	urlRaw, err := url.JoinPath(urlRelease, executableArchive.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = executableArchive, urlRaw

	return a, nil
}

func (m GitHub) ExecutableArchiveChecksums(v version.Version) (artifact.Remote[checksum.Executable], error) {
	var a artifact.Remote[checksum.Executable]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	checksumsExecutable, err := checksum.NewExecutable(v)
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	urlRaw, err := url.JoinPath(urlRelease, checksumsExecutable.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = checksumsExecutable, urlRaw

	return a, nil
}

func (m GitHub) SourceArchive(v version.Version) (artifact.Remote[source.Archive], error) {
	var a artifact.Remote[source.Archive]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v)
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	s := source.New(v)
	sourceArchive := source.Archive{Artifact: s}

	urlRaw, err := url.JoinPath(urlRelease, sourceArchive.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = sourceArchive, urlRaw

	return a, nil
}

func (m GitHub) SourceArchiveChecksums(v version.Version) (artifact.Remote[checksum.Source], error) {
	var a artifact.Remote[checksum.Source]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlRelease, err := urlGitHubRelease(v)
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	checksumsSource, err := checksum.NewSource(v)
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	urlRaw, err := url.JoinPath(urlRelease, checksumsSource.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = checksumsSource, urlRaw

	return a, nil
}

// Checks whether the version is broadly supported by the No network
// request is issued, but this does not guarantee the host has the version.
// To check whether the host has the version definitively via the network,
// use the 'CheckIfSupports' method.
func (m GitHub) Supports(v version.Version) bool {
	// GitHub only contains stable releases, starting with 'versionGitHubAssetSupport'.
	return v.IsStable() && v.CompareNormal(versionGitHubAssetSupport) >= 0
}

/* ----------------------- Function: urlGitHubRelease ----------------------- */

// Returns a URL to the version-specific release containing release assets.
func urlGitHubRelease(v version.Version) (string, error) {
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

	return url.JoinPath(gitHubAssetsURLBase, tag)
}