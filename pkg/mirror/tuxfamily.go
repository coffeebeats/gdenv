package mirror

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

const (
	tuxFamilyAssetsURLBase   = "https://downloads.tuxfamily.org/godotengine"
	tuxFamilyDirnameMono     = "mono"
	tuxFamilyDirnamePreAlpha = "pre-alpha"

	versionNormalV4 = "4.0.0"
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
	client *resty.Client
}

/* ------------------------- Function: NewTuxFamily ------------------------- */

func NewTuxFamily() TuxFamily {
	client := newClient()

	return TuxFamily{client}
}

/* ---------------------------- Method: Checksum ---------------------------- */

// Returns an 'Asset' to download the checksums file for the specified version
// from TuxFamily.
func (m *TuxFamily) Checksum(v godot.Version) (Asset, error) {
	var a Asset

	if !m.Supports(v) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, v.String())
	}

	urlVersionDir, err := urlTuxFamilyVersionDir(v)
	if err != nil {
		return a, err
	}

	urlRaw, err := url.JoinPath(urlVersionDir, filenameChecksums)
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
// from TuxFamily.
func (m *TuxFamily) Executable(ex godot.Executable) (Asset, error) {
	var a Asset

	if !m.Supports(ex.Version) {
		return a, fmt.Errorf("%w: '%s'", ErrInvalidSpecification, ex.Version.String())
	}

	name, err := ex.Name()
	if err != nil {
		return a, errors.Join(ErrInvalidSpecification, err)
	}

	filename := name + ".zip"

	urlVersionDir, err := urlTuxFamilyVersionDir(ex.Version)
	if err != nil {
		return a, err
	}

	urlRaw, err := url.JoinPath(urlVersionDir, filename)
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

	p = append(p, v.Normal())

	switch isMono, isStable := v.IsMono(), v.IsStable(); {
	case isMono:
		p = append(p, tuxFamilyDirnameMono)
	// For v4.0, the 'dev.*' labels are placed under the 'pre-alpha' directory.
	case v.Normal() == versionNormalV4 && reV4PreAlphaLabel.MatchString(v.Label()):
		p = append(p, tuxFamilyDirnamePreAlpha)
	case !isStable:
		p = append(p, v.Label())
	}

	urlVersionDir, err := url.JoinPath(tuxFamilyAssetsURLBase, p...)
	if err != nil {
		return "", errors.Join(ErrInvalidURL, err)
	}

	return urlVersionDir, nil
}
