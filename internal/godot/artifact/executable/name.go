package executable

import (
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                               Function: Name                               */
/* -------------------------------------------------------------------------- */

// Returns the name of the Godot executable, given the specified 'Version' and
// 'Platform'.
//
// NOTE: Godot names its executables in the format 'Godot_<VERSION>_<PLATFORM>',
// with Windows executables getting an extra '.exe' extension. Both the version
// and platform identifiers are version-specific, but the overall naming scheme
// has not changed (as of v4.2).
func Name(v version.Version, p platform.Platform) string {
	var name strings.Builder

	name.WriteString(namePrefix)
	name.WriteString(nameSeparator)

	name.WriteString(v.String())
	name.WriteString(nameSeparator)

	platformIdentifier, err := platform.Format(p, v)
	if err != nil {
		return ""
	}

	name.WriteString(platformIdentifier)

	if p.OS == platform.Windows {
		name.WriteString(".exe")
	}

	return name.String()
}
