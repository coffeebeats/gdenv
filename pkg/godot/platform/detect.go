package platform

import (
	"fmt"
	"os"
	"runtime"
)

const (
	EnvArch     = "GDENV_ARCH"
	EnvOS       = "GDENV_OS"
	EnvPlatform = "GDENV_PLATFORM"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Detect                              */
/* -------------------------------------------------------------------------- */

// Resolves the target platform by first checking environment variables and then
// falling back to the host platform.
func Detect() (Platform, error) {
	// First, check the full platform override.
	if platformRaw := os.Getenv(EnvPlatform); platformRaw != "" {
		p, err := Parse(platformRaw)
		if err != nil {
			return p, fmt.Errorf("%w: '%s'", err, platformRaw)
		}

		return p, nil
	}

	// Next, check the individual platform components for overrides and assemble
	// them into a 'Platform'.

	osRaw := os.Getenv(EnvOS)
	if osRaw == "" {
		osRaw = runtime.GOOS
	}

	o, err := ParseOS(osRaw)
	if err != nil {
		return Platform{}, fmt.Errorf("%w: '%s'", err, osRaw)
	}

	archRaw := os.Getenv(EnvArch)
	if archRaw == "" {
		archRaw = runtime.GOARCH
	}

	a, err := ParseArch(archRaw)
	if err != nil {
		return Platform{}, fmt.Errorf("%w: '%s'", err, archRaw)
	}

	return Platform{Arch: a, OS: o}, nil
}
