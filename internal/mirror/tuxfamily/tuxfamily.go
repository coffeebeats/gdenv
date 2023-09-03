package tuxfamily

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/pkg/godot"
)

const (
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
type TuxFamily struct {
	client client.Client
}

// Validate at compile-time that 'TuxFamily' implements 'Mirror'.
var _ mirror.Mirror = &TuxFamily{} //nolint:exhaustruct

/* ------------------------- Function: NewTuxFamily ------------------------- */

func NewTuxFamily() TuxFamily {
	client := client.New()

	return TuxFamily{client}
}

/* ---------------------------- Method: Checksum ---------------------------- */

// Returns an 'Asset' to download the checksums file for the specified version
// from TuxFamily.
func (m *TuxFamily) Checksum(v godot.Version) (mirror.Asset, error) {
	if !m.Supports(v) {
		return mirror.Asset{}, fmt.Errorf("%w: '%s'", mirror.ErrInvalidSpecification, v.String())
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return mirror.Asset{}, err
	}

	urlRaw, err := url.JoinPath(urlVersionDir, mirror.FilenameChecksums)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	return mirror.NewAsset(mirror.FilenameChecksums, urlRaw)
}

/* --------------------------- Method: Executable --------------------------- */

// Returns an 'Asset' to download a Godot executable for the specified version
// from TuxFamily.
func (m *TuxFamily) Executable(ex godot.Executable) (mirror.Asset, error) {
	if !m.Supports(ex.Version) {
		return mirror.Asset{}, fmt.Errorf("%w: '%s'", mirror.ErrInvalidSpecification, ex.Version.String())
	}

	name, err := ex.Name()
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidSpecification, err)
	}

	filename := name + ".zip"

	urlVersionDir, err := urlTuxFamilyVersionDir(ex.Version)
	if err != nil {
		return mirror.Asset{}, err
	}

	urlRaw, err := url.JoinPath(urlVersionDir, filename)
	if err != nil {
		return mirror.Asset{}, errors.Join(mirror.ErrInvalidURL, err)
	}

	return mirror.NewAsset(filename, urlRaw)
}

/* ------------------------------- Method: Has ------------------------------ */

// Returns whether the mirror supports the specified version. This does *NOT*
// guarantee that the mirror has the version.
func (m *TuxFamily) Supports(_ godot.Version) bool {
	return true
}

/* -------------------- Function: urlTuxFamilyVersionDir -------------------- */

// Returns a URL to the version-specific directory containing release assets.
//
// NOTE: Godot's TuxFamily directory structure is not straightforward. A
// route is built up in parts by replicating the directory structure. It's
// possible some edge cases are mishandled; please open an issue if one's found:
// https://github.com/coffeebeats/gdenv/issues/new?assignees=&labels=bug&projects=&template=%F0%9F%90%9B-bug-report.md
func urlTuxFamilyVersionDir(v godot.Version) (string, error) {
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
	case v.CompareNormal(godot.V4()) == 0 && reV4PreAlphaLabel.MatchString(v.Label()):
		p = append(p, tuxFamilyDirnamePreAlpha, strings.TrimPrefix(v.String(), godot.PrefixVersion))
	case !isStable:
		p = append(p, v.Label())
	}

	urlVersionDir, err := url.JoinPath(tuxFamilyAssetsURLBase, p...)
	if err != nil {
		return "", errors.Join(mirror.ErrInvalidURL, err)
	}

	return urlVersionDir, nil
}
