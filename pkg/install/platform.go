package install

import (
	"fmt"
	"os"
	"runtime"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
)

const (
	EnvGDEnvArch     = "GDENV_ARCH"
	EnvGDEnvOS       = "GDENV_OS"
	EnvGDEnvPlatform = "GDENV_PLATFORM"
)

/* -------------------------------------------------------------------------- */
/*                          Function: DetectPlatform                          */
/* -------------------------------------------------------------------------- */

// Resolves the target platform by first checking environment variables and then
// falling back to the host platform.
func DetectPlatform() (platform.Platform, error) {
	// First, check the full platform override.
	if platformRaw := os.Getenv(EnvGDEnvPlatform); platformRaw != "" {
		p, err := platform.Parse(platformRaw)
		if err != nil {
			return p, fmt.Errorf("%w: '%s'", err, platformRaw)
		}

		return p, nil
	}

	// Next, check the individual platform components for overrides and assemble
	// them into a 'Platform'.

	osRaw := os.Getenv(EnvGDEnvOS)
	if osRaw == "" {
		osRaw = runtime.GOOS
	}

	o, err := platform.ParseOS(osRaw)
	if err != nil {
		return platform.Platform{}, fmt.Errorf("%w: '%s'", err, osRaw)
	}

	archRaw := os.Getenv(EnvGDEnvArch)
	if archRaw == "" {
		archRaw = runtime.GOARCH
	}

	a, err := platform.ParseArch(archRaw)
	if err != nil {
		return platform.Platform{}, fmt.Errorf("%w: '%s'", err, archRaw)
	}

	return platform.Platform{Arch: a, OS: o}, nil
}
