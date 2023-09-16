package godot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const namePrefix = "Godot"
const nameSeparator = "_"

const (
	indexVersion    = 1
	indexPlatform   = nameSchemeParts - 1
	nameSchemeParts = 3 // Executables named in the format 'Godot_<VERSION>_<PLATFORM>'.
)

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
	Platform platform.Platform
	Version  version.Version
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

	platformIdentifier, err := platform.Format(e.Platform, e.Version)
	if err != nil {
		return "", err
	}

	name.WriteString(platformIdentifier)

	if e.Platform.OS == platform.Windows {
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
	if strings.HasPrefix(parts[indexPlatform], version.Mono) {
		parts[indexVersion] = parts[indexVersion] + nameSeparator + version.Mono
		parts[indexPlatform] = strings.TrimPrefix(parts[indexPlatform], version.Mono+nameSeparator)
	}

	v, err := version.Parse(parts[indexVersion])
	if err != nil {
		return executable, err
	}

	p, err := platform.Parse(parts[indexPlatform])
	if err != nil {
		return executable, err
	}

	executable.Platform = p
	executable.Version = v

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
