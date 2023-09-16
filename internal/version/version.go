package version

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

const (
	PrefixVersion              = "v"
	SeparatorBuildMetadata     = "+" // https://semver.org/#spec-item-10
	SeparatorPreReleaseVersion = "-" // https://semver.org/#spec-item-9

	LabelDefault = LabelStable
	LabelMono    = LabelStable + "_" + Mono
	LabelStable  = "stable"
	Mono         = "mono"
)

/* -------------------------------------------------------------------------- */
/*                        Functions: Version Constants                        */
/* -------------------------------------------------------------------------- */

// Returns a 'Version' struct for Godot v3.
func V3() Version {
	return Version{major: 3} //nolint:exhaustruct,gomnd
}

// Returns a 'Version' struct for Godot v4.
func V4() Version {
	return Version{major: 4} //nolint:exhaustruct,gomnd
}

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
		return LabelDefault
	}

	return v.label
}

/* ----------------------------- Method: IsMono ----------------------------- */

// Returns whether the version specifies a "mono" release (i.e. 'stable_mono').
func (v Version) IsMono() bool {
	return v.Label() == LabelMono
}

/* ---------------------------- Method: IsStable ---------------------------- */

// Returns whether the version specifies a "stable" release (e.g. 'stable' or
// 'stable_mono').
func (v Version) IsStable() bool {
	return v.Label() == LabelStable || v.Label() == LabelMono
}

/* ----------------------------- Method: Normal ----------------------------- */

// Returns the "normal version" format of the 'Version' (see
// https://semver.org/#spec-item-2).
func (v Version) Normal() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

/* -------------------------- Method: CompareNormal ------------------------- */

// Compares the "normal version" (see https://semver.org/#spec-item-2) to
// another 'Version' struct. The result will be '0' if 'v' == 'w', '-1' if
// 'v' < 'w', or '+1' if 'v' > 'w'.
func (v Version) CompareNormal(w Version) int {
	return semver.Compare(PrefixVersion+v.Normal(), PrefixVersion+w.Normal())
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
