package godot

import (
	"errors"
	"fmt"
	"strings"
)

const namePrefix = "Godot"
const nameSeparator = "_"

// Godot names its executables in the format 'Godot_<VERSION>_<PLATFORM>'.
const nameSchemeParts = 3

var (
	ErrInvalidName = errors.New("invalid name")
	ErrMissingName = errors.New("missing name")
)

/* -------------------------------------------------------------------------- */
/*                             Struct: Executable                             */
/* -------------------------------------------------------------------------- */

// A specification of a Godot executable (i.e. has a specific platform and
// version).
type Executable struct {
	Platform Platform
	Version  Version
}

/* ------------------------------ Method: Name ------------------------------ */

// Returns the name of the Godot executable, given the specified 'Version' and
// 'Platform'.
//
// NOTE: Godot names its executables in the format 'Godot_<VERSION>_<PLATFORM>',
// with Windows executables getting an extra '.exe' extension. Both the version
// and platform identifiers are version-specific, but the overall naming scheme
// has not changed (as of v4.2).
func (e Executable) Name() (string, error) {
	var name strings.Builder

	name.WriteString(namePrefix)
	name.WriteString(nameSeparator)

	name.WriteString(e.Version.String())
	name.WriteString(nameSeparator)

	platformIdentifier, err := FormatPlatform(e.Platform, e.Version)
	if err != nil {
		return "", err
	}

	name.WriteString(platformIdentifier)

	if e.Platform.OS == windows {
		name.WriteString(".exe")
	}

	return name.String(), nil
}

/* ----------------------------- Impl: Stringer ----------------------------- */

func (e Executable) String() string {
	name, err := e.Name()
	if err != nil {
		return ""
	}

	return name
}

/* -------------------------------------------------------------------------- */
/*                          Function: ParseExecutable                         */
/* -------------------------------------------------------------------------- */

// Parses an 'Executable' struct from the name of a Godot executable.
func ParseExecutable(input string) (Executable, error) {
	var executable Executable

	if input == "" {
		return executable, ErrMissingName
	}

	// In the event it's a Windows executable, trim the extension suffix. No
	// effect if it's a non-Windows build.
	input = strings.TrimSuffix(input, ".exe")

	// Try to split the input into 'Godot', '<VERSION>' and '<LABEL>'.
	parts := strings.SplitN(input, nameSeparator, nameSchemeParts)
	if len(parts) != nameSchemeParts {
		return executable, fmt.Errorf("%w: '%s'", ErrInvalidName, input)
	}

	// If the build is a "mono"-flavored build (i.e. the version label is
	// 'stable_mono'), then the third part will start with 'mono' instead of
	// the platform due to the "mono" version label containing the
	// 'nameSeparator' rune. Fix that here by removing the 'mono' prefix from
	// the platform and attaching it as a suffix to the version.
	if version, platform := 1, nameSchemeParts-1; strings.HasPrefix(parts[platform], mono) {
		parts[version] = parts[version] + nameSeparator + mono
		parts[platform] = strings.TrimPrefix(parts[platform], mono+nameSeparator)
	}

	version, err := ParseVersion(parts[1])
	if err != nil {
		return executable, err
	}

	platform, err := ParsePlatform(parts[2])
	if err != nil {
		return executable, err
	}

	executable.Platform = platform
	executable.Version = version

	return executable, nil
}

/* ---------------------- Function: MustParseExecutable --------------------- */

// Parses an 'Executable' struct from the name of a Godot executable or panics
// if it would fail.
func MustParseExecutable(input string) Executable {
	ex, err := ParseExecutable(input)
	if err != nil {
		panic(err)
	}

	return ex
}
