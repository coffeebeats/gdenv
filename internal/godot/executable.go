package godot

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

var (
	ErrInvalidVersion      = errors.New("godot: invalid version specification")
	ErrUnsupportedPlatform = errors.New("godot: unsupported platform")

	// Users can override platform detection by setting an environment variable.
	envVarPlatform = "GDENV_PLATFORM"
)

/* -------------------------- Function: Executable -------------------------- */

// Returns the name of the Godot executable, given the specified 'Version'.
func ExecutableName(v Version) (string, error) {
	if !v.IsValid() {
		return "", ErrInvalidVersion
	}

	var target string

	switch runtime.GOOS {
	case "darwin":
		target = "macos.universal"

	case "windows":
		switch runtime.GOARCH {
		case "386":
			target = "win32.exe"
		case "amd64":
			target = "win64.exe"
		}

	case "linux":
		switch runtime.GOARCH {
		case "386":
			target = "linux_x86_32"
		case "amd64":
			target = "linux_x86_64"
		}
	}

	// Override platform detection if the environment variable is set.
	if p := os.Getenv(envVarPlatform); p != "" {
		target = p
	}

	if target == "" {
		return "", fmt.Errorf("%w: %s/%s", ErrUnsupportedPlatform, runtime.GOOS, runtime.GOARCH)
	}

	// Set a default value.
	if v.Suffix == "" {
		v.Suffix = releaseLabelDefault
	}

	return fmt.Sprintf("Godot_%s_%s", v.String(), target), nil
}
