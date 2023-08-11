package godot

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

const releaseLabelDefault = "stable"

var (
	ErrInvalidInput = errors.New("invalid input string")
	ErrNoInput      = errors.New("no input provided")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Version                              */
/* -------------------------------------------------------------------------- */

// A struct containing a version specification for Godot.
type Version struct {
	Major  string
	Minor  string
	Patch  string
	Suffix string
}

/* ---------------------------- Method: Canonical --------------------------- */

// Returns a "canonical" representation of the 'Version', with all missing
// elements set to default values.
func (v *Version) Canonical() Version {
	major, minor, patch, suffix := v.Major, v.Minor, v.Patch, v.Suffix

	if major == "" {
		major = "0"
	}

	if minor == "" {
		minor = "0"
	}

	if patch == "" {
		patch = "0"
	}

	if suffix == "" {
		suffix = releaseLabelDefault
	}

	return Version{major, minor, patch, suffix}
}

/* ----------------------------- Method: IsValid ---------------------------- */

// Returns whether the 'Version' is well-specified.
func (v *Version) IsValid() bool {
	if v.Major == "" {
		return false
	}

	if v.Patch != "" && v.Minor == "" {
		return false
	}

	return true
}

/* ----------------------------- Impl: Stringer ----------------------------- */

// Returns an exact representation of the 'Version', if it's valid.
func (v *Version) String() string {
	// The 'Version' is not in a valid state.
	if !v.IsValid() {
		return ""
	}

	var s strings.Builder

	s.WriteRune('v')
	s.WriteString(v.Major)

	if v.Minor != "" {
		s.WriteRune('.')
		s.WriteString(v.Minor)
	}

	if v.Patch != "" {
		s.WriteRune('.')
		s.WriteString(v.Patch)
	}

	if v.Suffix != "" {
		s.WriteRune('-')
		s.WriteString(v.Suffix)
	}

	return s.String()
}

/* ------------------------- Function: ParseVersion ------------------------- */

// Parses a 'Version' struct from a non-canonical semantic version string.
func ParseVersion(s string) (Version, error) {
	var v Version

	if s == "" {
		return v, ErrNoInput
	}

	// Golang's 'semver' requires a 'v' prefix, but 'gdenv' doesn't.
	if !strings.HasPrefix(s, "v") {
		s = "v" + s
	}

	// Trim the version suffix, but store it for later.
	s, suffix, found := strings.Cut(s, "-")
	if (found && suffix == "") || !semver.IsValid(s) {
		return v, fmt.Errorf("%w: %s", ErrInvalidInput, s)
	}

	v.Suffix = suffix

	switch p := strings.Split(strings.TrimPrefix(s, "v"), "."); len(p) {
	case 3: //nolint:gomnd
		v.Patch = p[2]
		fallthrough // let 'Minor' and 'Major' be set
	case 2: //nolint:gomnd
		v.Minor = p[1]
		fallthrough // let 'Major' be set
	case 1:
		v.Major = p[0]
	default:
		return v, fmt.Errorf("%w: %s", ErrInvalidInput, s)
	}

	return v, nil
}
