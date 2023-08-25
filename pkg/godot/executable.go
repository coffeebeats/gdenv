package godot

import "strings"

const executableNamePrefix = "Godot"

/* -------------------------------------------------------------------------- */
/*                          Function: ExecutableName                          */
/* -------------------------------------------------------------------------- */

// Returns the name of the Godot executable, given the specified 'Version' and
// 'Platform'.
//
// NOTE: Godot names its executables in the format 'Godot_<VERSION>_<PLATFORM>',
// with Windows executables getting an extra '.exe' extension. Both the version
// and platform identifiers are version-specific, but the overall naming scheme
// has not changed (as of v4.2).
func ExecutableName(p Platform, v Version) (string, error) {
	var name strings.Builder

	name.WriteString(executableNamePrefix)
	name.WriteRune('_')

	name.WriteString(v.String())
	name.WriteRune('_')

	platformIdentifier, err := FormatPlatform(p, v)
	if err != nil {
		return "", err
	}

	name.WriteString(platformIdentifier)

	if p.os == windows {
		name.WriteString(".exe")
	}

	return name.String(), nil
}
