package godot

import (
	"errors"
	"fmt"
	"strings"
)

const namePrefix = "Godot"
const nameSeparator = '_'

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
	name.WriteRune(nameSeparator)

	name.WriteString(e.Version.String())
	name.WriteRune(nameSeparator)

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

	// Try to split the input into 'Godot_', '<VERSION>' and '<LABEL>'.
	parts := strings.SplitAfterN(input, string(nameSeparator), 2) //nolint:gomnd
	if len(parts) != 3 {                                          //nolint:gomnd
		return executable, fmt.Errorf("%w: '%s'", ErrInvalidName, input)
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
