package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

/* -------------------------------------------------------------------------- */
/*                           Function: ParseVersion                           */
/* -------------------------------------------------------------------------- */

// Parses a 'Version' struct from a semantic version string.
func ParseVersion(input string) (Version, error) {
	var version Version

	if input == "" {
		return version, ErrMissingVersion
	}

	// Normalize input by trimming excess space and using lowercase.
	input = strings.ToLower(strings.TrimSpace(input))

	// 'gdenv' and Semantic Versioning do not require a 'v' prefix (as of
	// version 2.0.0 - see semver.org/#semantic-versioning-200), but Golang's
	// 'semver' does.
	input = PrefixVersion + strings.TrimPrefix(input, PrefixVersion)

	// Trim the label off, but store it for later.
	input, label, found := strings.Cut(input, SeparatorPreReleaseVersion)
	if (found && label == "") || !semver.IsValid(input) {
		err := fmt.Errorf("%w: '%s'", ErrInvalidVersion, strings.TrimPrefix(input, PrefixVersion))
		return version, err
	}

	if label == labelDefault {
		label = ""
	}

	version.label = label

	// Trim build metadata - Godot does not use these (see https://semver.org/#spec-item-10).

	// NOTE: This step occurs *after* label extraction, so labels will keep any
	// build metadata suffixes. This future-proofs 'gdenv' if Godot changes its
	// build labeling practices. However, 'gdenv' doesn't support metadata
	// directly following the "normal version number".
	normalVersion, _, found := strings.Cut(input, SeparatorBuildMetadata)
	if found {
		return version, fmt.Errorf("%w: '%s'", ErrUnsupportedVersion, input)
	}

	parts, err := parseNormalVersion(normalVersion)
	if err != nil {
		return version, errors.Join(ErrInvalidVersion, err)
	}

	version.major, version.minor, version.patch = parts[0], parts[1], parts[2]

	return version, nil
}

/* ----------------------- Function: MustParseVersion ----------------------- */

// Parses a 'Version' struct from a semantic version string or panics if it
// would fail.
func MustParseVersion(input string) Version {
	v, err := ParseVersion(input)
	if err != nil {
		panic(err)
	}

	return v
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
		strings.Contains(input, SeparatorBuildMetadata) ||
		strings.Contains(input, SeparatorPreReleaseVersion) {
		panic(errNonNormalVersion)
	}

	// Remove the 'v' prefix to simplify version parsing below.
	input = strings.TrimPrefix(input, PrefixVersion)

	parts := strings.Split(input, ".")
	// NOTE: This should never occur for a 'semver'-validated string.
	if len(parts) < 1 || len(parts) > 3 {
		return out, fmt.Errorf("%w: '%s'", ErrUnrecognizedVersion, input)
	}

	for i, version := range parts { //nolint:varnamelen
		n, err := strconv.ParseUint(version, 10, 0)
		if err != nil {
			return out, fmt.Errorf("%w: '%s'", ErrInvalidVersionNumber, version)
		}

		out[i] = int(n)
	}

	return out, nil
}
