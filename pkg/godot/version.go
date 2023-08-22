package godot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

const labelDefault = "stable"

const prefixVersion = "v"
const separatorBuildMetadata = "+"     // https://semver.org/#spec-item-10
const separatorPreReleaseVersion = "-" // https://semver.org/#spec-item-9

var (
	ErrInvalidVersionInput      = errors.New("invalid input")
	ErrInvalidVersionNumber     = errors.New("invalid number")
	ErrMissingVersionInput      = errors.New("missing input")
	ErrUnrecognizedVersionInput = errors.New("unrecognized input")
	ErrUnsupportedVersion       = errors.New("unsupported version")

	errNonNormalVersion = errors.New("implementation error: found non-normal version")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Version                              */
/* -------------------------------------------------------------------------- */

// A struct containing a Godot version specification.
type Version struct {
	major, minor, patch int

	// Equivalent to "pre-release version" (see ttps://semver.org/#spec-item-9),
	// though Godot affixes "stable" to its stable releases. Note that an empty
	// 'Label' will be interpreted as a stable version.
	label string
}

/* ------------------------------ Method: Major ----------------------------- */

// Returns the major version component (see https://semver.org/#spec-item-4).
func (v Version) Major() int {
	return v.major
}

/* ------------------------------ Method: Minor ----------------------------- */

// Returns the minor version component (see https://semver.org/#spec-item-7).
func (v Version) Minor() int {
	return v.minor
}

/* ------------------------------ Method: Patch ----------------------------- */

// Returns the patch version component (see https://semver.org/#spec-item-6).
func (v Version) Patch() int {
	return v.patch
}

/* ------------------------------ Method: Label ----------------------------- */

// Returns the version label or a default if not defined.
func (v Version) Label() string {
	if v.label == "" {
		return labelDefault
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
func ParseVersion(input string) (Version, error) {
	var version Version

	if input == "" {
		return version, ErrMissingVersionInput
	}

	// 'gdenv' and Semantic Versioning do not require a 'v' prefix (as of
	// version 2.0.0 - see semver.org/#semantic-versioning-200), but Golang's
	// 'semver' does.
	input = prefixVersion + strings.TrimPrefix(input, prefixVersion)

	// Trim the label off, but store it for later.
	input, label, found := strings.Cut(input, separatorPreReleaseVersion)
	if (found && label == "") || !semver.IsValid(input) {
		return version, fmt.Errorf("%w: %s", ErrInvalidVersionInput, input)
	}

	version.label = label

	// Trim build metadata - Godot does not use these (see https://semver.org/#spec-item-10).

	// NOTE: This step occurs *after* label extraction, so labels will keep any
	// build metadata suffixes. This future-proofs 'gdenv' if Godot changes its
	// build labeling practices. However, 'gdenv' doesn't support metadata
	// directly following the "normal version number".
	normalVersion, _, found := strings.Cut(input, separatorBuildMetadata)
	if found {
		return version, fmt.Errorf("%w: %s", ErrUnsupportedVersion, input)
	}

	parts, err := parseNormalVersion(normalVersion)
	if err != nil {
		return version, errors.Join(ErrInvalidVersionInput, err)
	}

	version.major, version.minor, version.patch = parts[0], parts[1], parts[2]

	return version, nil
}

/* ---------------------- Function: parseNormalVersion ---------------------- */

// Parses the "normal version" (see https://semver.org/#spec-item-2) from a
// 'semver'-validated version string.
//
// NOTE: This implementation requires that there are *no* build or version
// specifiers.
func parseNormalVersion(input string) ([3]int, error) {
	out := [3]int{0, 0, 0}

	if !semver.IsValid(input) ||
		strings.Contains(input, separatorBuildMetadata) ||
		strings.Contains(input, separatorPreReleaseVersion) {
		panic(errNonNormalVersion)
	}

	// Remove the 'v' prefix to simplify version parsing below.
	input = strings.TrimPrefix(input, prefixVersion)

	parts := strings.Split(input, ".")
	// NOTE: This should never occur for a 'semver'-validated string.
	if len(parts) < 1 || len(parts) > 3 {
		return out, fmt.Errorf("%w: %s", ErrUnrecognizedVersionInput, input)
	}

	for i, version := range parts { //nolint:varnamelen
		n, err := strconv.ParseUint(version, 10, 0)
		if err != nil {
			return out, fmt.Errorf("%w: %s", ErrInvalidVersionNumber, version)
		}

		out[i] = int(n)
	}

	return out, nil
}
