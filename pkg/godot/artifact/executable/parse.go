package executable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

const (
	indexVersion  = 1
	indexPlatform = nameSchemeParts - 1

	nameSchemeParts = 3 // Executables named in the format 'Godot_<VERSION>_<PLATFORM>'.
)

var (
	ErrInvalidName = errors.New("invalid name")
	ErrMissingName = errors.New("missing name")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Parse                              */
/* -------------------------------------------------------------------------- */

// Parses an 'Executable' struct from the name of a Godot executable.
func Parse(input string) (Executable, error) {
	var ex Executable

	if input == "" {
		return ex, ErrMissingName
	}

	// In the event it's a Windows executable, trim the extension suffix. No
	// effect if it's a non-Windows build.
	input = strings.TrimSuffix(input, ".exe")

	// Try to split the input into 'Godot', '<VERSION>' and '<LABEL>'.
	parts := strings.SplitN(input, nameSeparator, nameSchemeParts)
	if len(parts) != nameSchemeParts {
		return ex, fmt.Errorf("%w: '%s'", ErrInvalidName, input)
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
		return ex, err
	}

	p, err := platform.Parse(parts[indexPlatform])
	if err != nil {
		return ex, err
	}

	ex.platform = p
	ex.version = v

	return ex, nil
}

/* -------------------------------------------------------------------------- */
/*                             Function: MustParse                            */
/* -------------------------------------------------------------------------- */

// Parses an 'Executable' struct from the name of a Godot executable or panics
// if it would fail.
func MustParse(input string) Executable {
	ex, err := Parse(input)
	if err != nil {
		panic(err)
	}

	return ex
}
