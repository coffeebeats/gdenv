package mirror

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	tuxfamilyDownloadsDomain = "downloads.tuxfamily.org"
	tuxFamilyAssetsURLBase   = "https://downloads.tuxfamily.org/godotengine"
	tuxFamilyDirnameMono     = "mono"
	tuxFamilyDirnamePreAlpha = "pre-alpha"
)

var (
	// This expression matches all Godot v4.0 "pre-alpha" versions which use a
	// release label similar to 'dev.20211015'. This expressions has been tested
	// manually.
	reV4PreAlphaLabel = regexp.MustCompile(`^dev\.[0-9]{8}$`)
)

/* -------------------------------------------------------------------------- */
/*                              Struct: TuxFamily                             */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via the Godot TuxFamily host.
type TuxFamily[T artifact.Versioned] struct{}

// Validate at compile-time that 'TuxFamily' implements 'Mirror' interfaces.
var _ Hoster = &TuxFamily[artifact.Versioned]{}
var _ Remoter[executable.Archive] = &TuxFamily[executable.Archive]{}
var _ Remoter[source.Archive] = &TuxFamily[source.Archive]{}

/* ------------------------------ Impl: Hoster ------------------------------ */

// Hosts returns the host URLs at which artifacts are hosted.
func (m TuxFamily[T]) Hosts() []string {
	return []string{gitHubContentDomain}
}

/* ------------------------------ Impl: Remoter ----------------------------- */

// Remote returns an 'artifact.Remote' wrapper around a specified artifact. The
// remote wrapper contains the URL at which the artifact can be downloaded.
func (m TuxFamily[T]) Remote(a T) (artifact.Remote[T], error) {
	var remote artifact.Remote[T]

	urlVersionDir, err := urlTuxFamilyVersionDir(a.Version())
	if err != nil {
		return remote, err
	}

	urlParsed, err := client.ParseURL(urlVersionDir, a.Name())
	if err != nil {
		return remote, errors.Join(ErrInvalidURL, err)
	}

	remote.Artifact, remote.URL = a, urlParsed

	return remote, nil
}

/* ------------------------------ Impl: Mirror ------------------------------ */

// Name returns the display name of the mirror.
func (m TuxFamily[T]) Name() string {
	return "TuxFamily (downloads.tuxfamily.org/godotengine)"
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
