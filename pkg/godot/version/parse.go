package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

var (
	ErrInvalid       = errors.New("invalid version")
	ErrInvalidNumber = errors.New("invalid version number")
	ErrMissing       = errors.New("missing version")
	ErrUnrecognized  = errors.New("unrecognized version")
	ErrUnsupported   = errors.New("unsupported version")

	errNonNormalVersion = errors.New("implementation error: found non-normal version")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Parse                              */
/* -------------------------------------------------------------------------- */

// Parses a 'Version' struct from a semantic version string.
func Parse(input string) (Version, error) {
	var version Version

	if input == "" {
		return version, ErrMissing
	}

	// Normalize input by trimming excess space and using lowercase.
	input = strings.ToLower(strings.TrimSpace(input))

	// 'gdenv' and Semantic Versioning do not require a 'v' prefix (as of
	// version 2.0.0 - see semver.org/#semantic-versioning-200), but Golang's
	// 'semver' does.
	input = Prefix + strings.TrimPrefix(input, Prefix)

	// Trim the label off, but store it for later.
	input, label, found := strings.Cut(input, SeparatorPreReleaseVersion)
	if (found && label == "") || !semver.IsValid(input) {
		err := fmt.Errorf("%w: '%s'", ErrInvalid, strings.TrimPrefix(input, Prefix))

		return version, err
	}

	if label == LabelDefault {
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
		return version, fmt.Errorf("%w: '%s'", ErrUnsupported, input)
	}

	parts, err := parseNormalVersion(normalVersion)
	if err != nil {
		return version, errors.Join(ErrInvalid, err)
	}

	version.major, version.minor, version.patch = parts[0], parts[1], parts[2]

	return version, nil
}

/* --------------------------- Function: MustParse -------------------------- */

// Parses a 'Version' struct from a semantic version string or panics if it
// would fail.
func MustParse(input string) Version {
	v, err := Parse(input)
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
	input = strings.TrimPrefix(input, Prefix)

	parts := strings.Split(input, ".")
	// NOTE: This should never occur for a 'semver'-validated string.
	if len(parts) < 1 || len(parts) > 3 {
		return out, fmt.Errorf("%w: '%s'", ErrUnrecognized, input)
	}

	for i, version := range parts { //nolint:varnamelen
		n, err := strconv.ParseUint(version, 10, 32)
		if err != nil {
			return out, fmt.Errorf("%w: '%s'", ErrInvalidNumber, version)
		}

		out[i] = int(n)
	}

	return out, nil
}
