package godot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

const labelStable = "stable"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidNumber      = errors.New("invalid number")
	ErrMissingInput       = errors.New("missing input")
	ErrUnsupportedVersion = errors.New("unsupported version")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Version                              */
/* -------------------------------------------------------------------------- */

// A struct containing a Godot version specification.
type Version struct {
	major, minor, patch int

	// Equivalent to "pre-release version" (see semver.org), though Godot
	// affixes "stable" to its stable releases. Note that an empty 'Label'
	// will be interpreted as a stable version.
	label string
}

/* ------------------------------ Method: Major ----------------------------- */

// Returns the major version component as an 'int' (see https://semver.org/#spec-item-4).
func (v Version) Major() int {
	return v.major
}

/* ------------------------------ Method: Minor ----------------------------- */

// Returns the minor version component as an 'int' (see https://semver.org/#spec-item-7).
func (v Version) Minor() int {
	return v.minor
}

/* ------------------------------ Method: Patch ----------------------------- */

// Returns the patch version component as an 'int' (see https://semver.org/#spec-item-6).
func (v Version) Patch() int {
	return v.patch
}

/* ------------------------------ Method: Label ----------------------------- */

// Returns the version label associated with the Godot version specification.
func (v Version) Label() string {
	if v.label == "" {
		return labelStable
	}

	return v.label
}

/* ----------------------------- Impl: Stringer ----------------------------- */

func (v Version) String() string {
	var out strings.Builder

	out.WriteRune('v')

	out.WriteString(strconv.Itoa(v.Major()))
	out.WriteRune('.')

	out.WriteString(strconv.Itoa(v.Minor()))

	// Godot never includes a trailing '.0' for patch versions (as of v4.2).
	if patch := v.Patch(); patch != 0 {
		out.WriteRune('.')
		out.WriteString(strconv.Itoa(patch))
	}

	out.WriteRune('-')
	out.WriteString(v.Label())

	return out.String()
}

/* -------------------------------------------------------------------------- */
/*                           Function: ParseVersion                           */
/* -------------------------------------------------------------------------- */

// Parses a 'Version' struct from a semantic version string.
func ParseVersion(input string) (Version, error) { //nolint:funlen
	var out Version

	if input == "" {
		return out, ErrMissingInput
	}

	// Golang's 'semver' requires a 'v' prefix, but 'gdenv' and Semantic
	// Versioning do not (as of version 2.0.0 - see semver.org/#semantic-versioning-200).
	if !strings.HasPrefix(input, "v") {
		input = "v" + input
	}

	// Trim the label off, but store it for later.
	version, label, found := strings.Cut(input, "-")
	if (found && label == "") || !semver.IsValid(version) {
		return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
	}

	out.label = label

	// Trim build metadata - Godot does not use these (see https://semver.org/#spec-item-10).

	// NOTE: This step occurs *after* label extraction, so labels will keep any
	// build metadata suffixes. This future-proofs 'gdenv' if Godot changes its build
	// labeling practices. However, 'gdenv' doesn't support metadata directly
	// following the "normal version number".
	version, _, found = strings.Cut(version, "+")
	if found {
		return out, fmt.Errorf("%w: %s", ErrUnsupportedVersion, version)
	}

	switch parts := strings.Split(version, "."); len(parts) {
	case 3: //nolint:gomnd
		n, err := parseNumber(parts[2])
		if err != nil {
			return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
		}

		out.patch = n

		fallthrough // let 'Minor' and 'Major' be set
	case 2: //nolint:gomnd
		n, err := parseNumber(parts[1])
		if err != nil {
			return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
		}

		out.minor = n

		fallthrough // let 'Major' be set
	case 1:
		n, err := parseNumber(strings.TrimPrefix(parts[0], "v"))
		if err != nil {
			return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
		}

		out.major = n
	default:
		return out, fmt.Errorf("%w: %s", ErrInvalidInput, version)
	}

	return out, nil
}

/* -------------------------- Function: parseNumber ------------------------- */

// Parses an unsigned integer from a string, but returns an 'int'. This is a
// convenience function when parsing Semantic Versioning components, ensuring
// integers are greater than '0'.
func parseNumber(s string) (int, error) {
	n, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidNumber, s)
	}

	return int(n), nil
}
