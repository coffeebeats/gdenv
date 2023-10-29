package mirror

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	tuxfamilyDownloadsDomain = "downloads.tuxfamily.org"
	tuxFamilyAssetsURLBase   = "https://downloads.tuxfamily.org/godotengine"
	tuxFamilyDirnameMono     = "mono"
	tuxFamilyDirnamePreAlpha = "pre-alpha"
)

var (
	versionTuxFamilyMinSupported = version.MustParse("v1.1") //nolint:gochecknoglobals

	// This expression matches all Godot v4.0 "pre-alpha" versions which use a
	// release label similar to 'dev.20211015'. This expressions has been tested
	// manually.
	reV4PreAlphaLabel = regexp.MustCompile(`^dev\.[0-9]{8}$`)
)

/* -------------------------------------------------------------------------- */
/*                              Struct: TuxFamily                             */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via the Godot TuxFamily host.
type TuxFamily struct{}

// Validate at compile-time that 'TuxFamily' implements 'Mirror' interfaces.
var _ Mirror = &TuxFamily{}
var _ Executable = &TuxFamily{}
var _ Source = &TuxFamily{}

/* ------------------------------ Impl: Mirror ------------------------------ */

// Returns a new 'client.Client' for downloading artifacts from the mirror.
func (m TuxFamily) Domains() []string {
	return nil
}

// Checks whether the version is broadly supported by the mirror. No network
// request is issued, but this does not guarantee the host has the version.
// To check whether the host has the version definitively via the network,
// use the 'checkIfExists' method.
func (m TuxFamily) Supports(v version.Version) bool {
	// TuxFamily seems to contain all published releases.
	return v.CompareNormal(versionTuxFamilyMinSupported) >= 0
}

/* ---------------------------- Impl: Executable ---------------------------- */

func (m TuxFamily) ExecutableArchive(
	v version.Version,
	p platform.Platform,
) (artifact.Remote[executable.Archive], error) {
	var a artifact.Remote[executable.Archive]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v)
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return a, err
	}

	executableArchive := executable.Archive{Artifact: executable.New(v, p)}

	urlParsed, err := client.ParseURL(urlVersionDir, executableArchive.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = executableArchive, urlParsed

	return a, nil
}

func (m TuxFamily) ExecutableArchiveChecksums(v version.Version) (artifact.Remote[checksum.Executable], error) {
	var a artifact.Remote[checksum.Executable]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	checksumsExecutable, err := checksum.NewExecutable(v)
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return a, err
	}

	urlParsed, err := client.ParseURL(urlVersionDir, checksumsExecutable.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = checksumsExecutable, urlParsed

	return a, nil
}

/* ------------------------------ Impl: Source ------------------------------ */

func (m TuxFamily) SourceArchive(v version.Version) (artifact.Remote[source.Archive], error) {
	var a artifact.Remote[source.Archive]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return a, err
	}

	s := source.New(v)
	sourceArchive := source.Archive{Artifact: s}

	urlParsed, err := client.ParseURL(urlVersionDir, sourceArchive.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = sourceArchive, urlParsed

	return a, nil
}

func (m TuxFamily) SourceArchiveChecksums(v version.Version) (artifact.Remote[checksum.Source], error) {
	var a artifact.Remote[checksum.Source]

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	checksumsSource, err := checksum.NewSource(v)
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return a, err
	}

	urlParsed, err := client.ParseURL(urlVersionDir, checksumsSource.Name())
	if err != nil {
		return a, errors.Join(ErrInvalidURL, err)
	}

	a.Artifact, a.URL = checksumsSource, urlParsed

	return a, nil
}

/* -------------------- Function: urlTuxFamilyVersionDir -------------------- */

// Returns a URL to the version-specific directory containing release assets.
//
// NOTE: Godot's TuxFamily directory structure is not straightforward. A
// route is built up in parts by replicating the directory structure. It's
// possible some edge cases are mishandled; please open an issue if one's found:
// https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=%F0%9F%90%9B-bug-report.md
func urlTuxFamilyVersionDir(v version.Version) (string, error) {
	p := make([]string, 0)

	// The first directory will be the "normal version", but a patch version of
	// '0' will be dropped.
	var normal string

	switch v.Patch() {
	case 0:
		normal = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
	default:
		normal = v.Normal()
	}

	p = append(p, normal)

	// If the build is a "stable", non-"mono" flavor, then the assets will be in
	// the version directory. Otherwise, the assets will be in one or more sub-
	// directories corresponding to build labels.
	switch isMono, isStable := v.IsMono(), v.IsStable(); {
	case isMono:
		p = append(p, tuxFamilyDirnameMono)
	// For v4.0, the 'dev.*' labels are placed under label subdirectories,
	// themselves under the 'pre-alpha' directory.
	case v.CompareNormal(version.Godot4()) == 0 && reV4PreAlphaLabel.MatchString(v.Label()):
		p = append(p, tuxFamilyDirnamePreAlpha, strings.TrimPrefix(v.String(), version.Prefix))
	case !isStable:
		p = append(p, v.Label())
	}

	urlVersionDir, err := url.JoinPath(tuxFamilyAssetsURLBase, p...)
	if err != nil {
		return "", errors.Join(ErrInvalidURL, err)
	}

	return urlVersionDir, nil
}
