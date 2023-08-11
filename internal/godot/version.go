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

	var out strings.Builder

	out.WriteRune('v')
	out.WriteString(v.Major)

	if v.Minor != "" {
		out.WriteRune('.')
		out.WriteString(v.Minor)
	}

	if v.Patch != "" {
		out.WriteRune('.')
		out.WriteString(v.Patch)
	}

	if v.Suffix != "" {
		out.WriteRune('-')
		out.WriteString(v.Suffix)
	}

	return out.String()
}

/* ------------------------- Function: ParseVersion ------------------------- */

// Parses a 'Version' struct from a non-canonical semantic version string.
func ParseVersion(version string) (Version, error) {
	var out Version

	if version == "" {
		return out, ErrNoInput
	}

	// Golang's 'semver' requires a 'v' prefix, but 'gdenv' doesn't.
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// Trim a valid build label suffix; Godot does not use these.
	version, _, _ = strings.Cut(version, "+")

	// Trim the version suffix, but store it for later.
	version, suffix, found := strings.Cut(version, "-")
	if (found && suffix == "") || !semver.IsValid(version) {
		return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
	}

	out.Suffix = suffix

	switch parts := strings.Split(strings.TrimPrefix(version, "v"), "."); len(parts) {
	case 3: //nolint:gomnd
		out.Patch = parts[2]
		fallthrough // let 'Minor' and 'Major' be set
	case 2: //nolint:gomnd
		out.Minor = parts[1]
		fallthrough // let 'Major' be set
	case 1:
		out.Major = parts[0]
	default:
		return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
	}

	return out, nil
}
