package godot

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

var (
	ErrMissingArch          = errors.New("missing architecture")
	ErrMissingOS            = errors.New("missing OS")
	ErrMissingPlatform      = errors.New("missing platform")
	ErrUnrecognizedArch     = errors.New("unrecognized architecture")
	ErrUnrecognizedOS       = errors.New("unrecognized OS")
	ErrUnrecognizedPlatform = errors.New("unrecognized platform")
	ErrUnsupportedArch      = errors.New("unsupported architecture")
	ErrUnsupportedOS        = errors.New("unsupported OS")

	// This expression matches all Godot v4.0 pre-release versions which utilize
	// a 'osx.universal' platform label. These include 'alpha1' - 'alpha12' and
	// all of the 'dev.*' pre-alpha versions. This expressions has been tested
	// manually and some unit tests validate this as well.
	reV4LabelsWithOSXUniversal = regexp.MustCompile(`^(alpha([1-9]|1[0-2])|(dev\.[0-9]{8}))$`)
)

/* -------------------------------------------------------------------------- */
/*                                  Enum: OS                                  */
/* -------------------------------------------------------------------------- */

// Operating systems which the Godot project provides prebuilt editor binaries
// for.
type OS int

const (
	linux OS = iota + 1
	macOS
	windows
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
		return 0, ErrMissingOS

	case "darwin", "macos", "osx":
		return macOS, nil

	case "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
		return linux, nil

	case "win", "windows":
		return windows, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedOS, input)
	}
}

/* -------------------------- Function: MustParseOS ------------------------- */

// Parses an input string as an operating system specification but panics if it
// would fail.
func MustParseOS(input string) OS {
	os, err := ParseOS(input)
	if err != nil {
		panic(err)
	}

	return os
}

/* -------------------------------------------------------------------------- */
/*                                 Enum: Arch                                 */
/* -------------------------------------------------------------------------- */

// CPU architectures which the Godot project provides prebuilt editor binaries
// for.
type Arch int

const (
	amd64 Arch = iota + 1
	arm64
	i386
	universal
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
		return 0, ErrMissingArch

	case "386", "i386", "x86", "x86_32":
		return i386, nil

	case "amd64", "x86_64", "x86-64":
		return amd64, nil

	case "arm64", "arm64be":
		return arm64, nil

	case "fat", "universal":
		return universal, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedArch, input)
	}
}

/* ------------------------- Function: MustParseArch ------------------------ */

// Parses an input string as a CPU architecture specification but panics if it
// would fail.
func MustParseArch(input string) Arch {
	arch, err := ParseArch(input)
	if err != nil {
		panic(err)
	}

	return arch
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

	os, err := ParseOS(runtime.GOOS)
	if err != nil {
		return platform, fmt.Errorf("%w: '%s'", err, runtime.GOOS)
	}

	arch, err := ParseArch(runtime.GOARCH)
	if err != nil {
		return platform, fmt.Errorf("%w: '%s'", err, runtime.GOOS)
	}

	platform.Arch, platform.OS = arch, os

	return platform, nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: ParsePlatform                          */
/* -------------------------------------------------------------------------- */

// Parses a 'Platform' struct from a platform identifier. There are potentially
// multiple valid identifiers for any given platform due to schema differences
// across Godot versions.
func ParsePlatform(input string) (Platform, error) {
	if input == "" {
		return Platform{}, ErrMissingPlatform
	}

	switch strings.ToLower(strings.TrimSpace(input)) {
	// Linux
	case "x11.32", "linux.x86_32":
		return Platform{i386, linux}, nil
	case "x11.64", "linux.x86_64":
		return Platform{amd64, linux}, nil

	// MacOS - Note that the supported architectures between 'osx.fat' and
	// 'osx.universal' are *NOT* the same. It's important to maintain the
	// 'Version' alongside this result so that the architectures can be
	// correctly determined.
	case "osx.64":
		return Platform{amd64, macOS}, nil
	case "macos.universal", "osx.fat", "osx.universal":
		return Platform{universal, macOS}, nil

	// Windows
	case "win32":
		return Platform{i386, windows}, nil
	case "win64":
		return Platform{amd64, windows}, nil

	default:
		return Platform{}, fmt.Errorf("%w: '%s'", ErrUnrecognizedPlatform, input)
	}
}

/* ----------------------- Function: MustParsePlatform ---------------------- */

// Parses an input string as a 'Platform' specification but panics if it would
// fail.
func MustParsePlatform(input string) Platform {
	platform, err := ParsePlatform(input)
	if err != nil {
		panic(err)
	}

	return platform
}

/* -------------------------------------------------------------------------- */
/*                          Function: FormatPlatform                          */
/* -------------------------------------------------------------------------- */

// Formats a 'Platform' specification into a platform string found in Godot
// executable names.
//
// NOTE: This method delegates to rather complex OS-specific methods. It would
// be great if there were a better way to organize this.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some platform identifiers are missing or incorrect:
// github.com/coffeebeats/gdenv/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func FormatPlatform(p Platform, v Version) (string, error) {
	switch p.OS {
	case linux:
		return formatLinuxPlatform(p.Arch, v)
	case macOS:
		return formatMacOSPlatform(p.Arch, v)
	case windows:
		return formatWindowsPlatform(p.Arch, v)

	case 0:
		return "", ErrMissingOS
	default:
		return "", ErrUnrecognizedOS
	}
}

/* ---------------------- Function: formatLinuxPlatform --------------------- */

// Given an architecture, returns the Linux platform identifier used by Godot
// executable names.
func formatLinuxPlatform(a Arch, v Version) (string, error) { //nolint:cyclop
	if a == 0 {
		return "", ErrMissingArch
	}

	var p string

	switch {
	// Godot v1-v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", ErrUnsupportedVersion, v)
	// Godot v3
	case v.major < 4: //nolint:gomnd
		// 'linux_headless.64' and 'linux_server.64' flavors introduced in v3.1
		// are not supported.
		switch a {
		case i386:
			p = "x11.32"
		case amd64:
			p = "x11.64"

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// Godot v4+
	default:
		switch a {
		case i386:
			p = "linux.x86_32"
		case amd64:
			p = "linux.x86_64"

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}

	// For some reason, "mono"-flavored builds have the '.' rune replaced by a
	// '_' rune.
	if v.IsMono() {
		p = strings.ReplaceAll(p, ".", "_")
	}

	return p, nil
}

/* ---------------------- Function: formatMacOSPlatform --------------------- */

// Given an architecture, returns the macOS platform identifier used by Godot
// executable names.
//
// NOTE: This is rather convoluted; consider a better way of organizing this
// logic.
func formatMacOSPlatform(a Arch, v Version) (string, error) { //nolint:cyclop
	if a == 0 {
		return "", ErrMissingArch
	}

	switch {
	// Godot v1 - v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", ErrUnsupportedVersion, v)
	// Godot v3.0 - v3.0.6
	case v.major == 3 && v.minor < 1:
		switch a {
		case i386, amd64:
			return "osx.fat", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// Godot v3.1 - v3.2.4-beta2
	// NOTE: Because v3.2.4 labels are only "beta" and "rc" *and* "beta"
	// versions do not exceed 6, lexicographic  sorting works.
	case v.major == 3 && v.minor <= 2 && (v.patch < 4 || v.patch == 4 && v.label <= "beta2"):
		switch a {
		case amd64:
			return "osx.64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// Godot v3.2.4-beta3 - v4.0-alpha12
	case v.CompareNormal(godotVersion4()) < 0 ||
		(v.CompareNormal(godotVersion4()) == 0 && reV4LabelsWithOSXUniversal.MatchString(v.label)):
		switch a {
		case amd64, arm64:
			return "osx.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// Godot v4.0-alpha13+
	default:
		switch a {
		case amd64, arm64:
			return "macos.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}
}

/* --------------------- Function: formatWindowsPlatform -------------------- */

// Given an architecture, returns the Windows platform identifier used by Godot
// executable names.
func formatWindowsPlatform(a Arch, v Version) (string, error) {
	if a == 0 {
		return "", ErrMissingArch
	}

	switch {
	// Godot v1-v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", ErrUnsupportedVersion, v)
	// Godot v3+
	default:
		switch a {
		case i386:
			return "win32", nil
		case amd64:
			return "win64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: godotVersion4                          */
/* -------------------------------------------------------------------------- */

func godotVersion4() Version {
	return Version{major: 4} //nolint:exhaustruct,gomnd
}
