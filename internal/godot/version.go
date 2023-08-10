package godot

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

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

/* ----------------------------- Impl: Stringer ----------------------------- */

func (v *Version) String() string {
	major := v.Major
	if major == "" {
		major = "0"
	}

	minor := v.Minor
	if minor == "" {
		minor = "0"
	}

	patch := v.Patch
	if patch == "" {
		patch = "0"
	}

	suffix := v.Suffix
	if suffix == "" {
		suffix = "stable"
	}

	return fmt.Sprintf("v%s.%s.%s-%s", major, minor, patch, suffix)
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

	version, suffix, found := strings.Cut(s, "-")
	if !found {
		suffix = "stable"
	}

	if suffix == "" {
		return v, fmt.Errorf("%w: %s", ErrInvalidInput, s)
	}

	version = strings.TrimPrefix(semver.Canonical(version), "v")
	if len(version) == 0 {
		return v, fmt.Errorf("%w: %s", ErrInvalidInput, s)
	}

	parts := strings.Split(version, ".")

	v.Major, v.Minor, v.Patch = parts[0], parts[1], parts[2]
	v.Suffix = suffix

	return v, nil
}
