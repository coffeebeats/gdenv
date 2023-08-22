package godot

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrMissingTargetInput      = errors.New("missing input")
	ErrUnrecognizedTargetInput = errors.New("unrecognized input")
)

/* -------------------------------------------------------------------------- */
/*                                  Enum: OS                                  */
/* -------------------------------------------------------------------------- */

// Operating systems which the Godot project provides prebuilt binaries for.
type OS int

const (
	Linux OS = iota + 1
	MacOS
	Windows
)

/* ---------------------------- Function: ParseOS --------------------------- */

// Parses an input string as an operating system specification. Typically this
// will rely on the 'GOOS' value, but users may override that setting via
// an environment variable. As such, this function recognizes some values in
// addition to what Go defines.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdenv/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseOS(input string) (OS, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingTargetInput

	case "darwin", "macos", "osx":
		return MacOS, nil

	case "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
		return Linux, nil

	case "win", "windows":
		return Windows, nil

	default:
		return 0, fmt.Errorf("%w: %s", ErrUnrecognizedTargetInput, input)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 Enum: Arch                                 */
/* -------------------------------------------------------------------------- */

// CPU architectures which the Godot project provides prebuilt binaries for.
type Arch int

const (
	Amd64 Arch = iota + 1
	Arm64
	I386
	Universal
)

/* --------------------------- Function: ParseArch -------------------------- */

// Parses an input string as a CPU architecture specification. Typically this
// will rely on the 'os.GOARCH' value, but users may override that setting via
// an environment variable. As such, this function recognizes some values in
// addition to what Go defines.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdenv/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseArch(input string) (Arch, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingTargetInput

	case "386", "i386", "x86", "x86_32":
		return I386, nil

	case "amd64", "x86_64", "x86-64":
		return Amd64, nil

	case "arm64", "arm64be":
		return Arm64, nil

	case "fat", "universal":
		return Universal, nil

	default:
		return 0, fmt.Errorf("%w: %s", ErrUnrecognizedTargetInput, input)
	}
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Platform                              */
/* -------------------------------------------------------------------------- */

// A platform specification representing a target to run the Godot editor on.
type Platform struct {
	Arch Arch
	OS   OS
}

/* ------------------------- Function: HostPlatform ------------------------- */

// Returns a 'Platform' struct pertaining to the host machine, if recognized.
func HostPlatform() (Platform, error) {
	var platform Platform

	oS, err := ParseOS(runtime.GOOS)
	if err != nil {
		return platform, fmt.Errorf("%w: %s", err, runtime.GOOS)
	}

	arch, err := ParseArch(runtime.GOARCH)
	if err != nil {
		return platform, fmt.Errorf("%w: %s", err, runtime.GOOS)
	}

	platform.Arch, platform.OS = arch, oS

	return platform, nil
}
