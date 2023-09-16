package platform

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot"
)

var (
	ErrMissingPlatform      = errors.New("missing platform")
	ErrUnrecognizedPlatform = errors.New("unrecognized platform")

	// This expression matches all Godot v4.0 macOS pre-release versions which
	// utilize a 'osx.universal' platform label. These include 'alpha1' -
	// 'alpha12' and all of the 'dev.*' pre-alpha versions. This expressions has
	// been tested manually and some unit tests validate this as well.
	reV4MacOSLabelsWithOSXUniversal = regexp.MustCompile(`^(alpha([1-9]|1[0-2])|(dev\.[0-9]{8}))$`)

	// This expression matches all Godot v4.0 Linux pre-release versions which
	// utilize a platform label. These include 'alpha1' - 'alpha14' and all of
	// the 'dev.*' pre-alpha versions. This expressions has been tested manually
	// and some unit tests validate this as well.
	reV4LinuxLabelsWithoutX86 = regexp.MustCompile(`^(alpha([1-9]|1[0-4])|(dev\.[0-9]{8}))$`)
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Platform                              */
/* -------------------------------------------------------------------------- */

// A platform specification representing a target to run the Godot editor on.
type Platform struct {
	Arch Arch
	OS   OS
}

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Platform' struct from a valid 'OS' and 'Arch'.
func New(os OS, arch Arch) (Platform, error) {
	var platform Platform

	switch os {
	case Linux, MacOS, Windows:

	case 0:
		return platform, ErrMissingOS
	default:
		return platform, fmt.Errorf("%w: '%d'", ErrUnrecognizedOS, os)
	}

	switch arch {
	case Amd64, Arm64, I386, Universal:

	case 0:
		return platform, ErrMissingArch
	default:
		return platform, fmt.Errorf("%w: '%d'", ErrUnrecognizedArch, arch)
	}

	platform.Arch = arch
	platform.OS = os

	return platform, nil
}

/* -------------------------------------------------------------------------- */
/*                               Function: Parse                              */
/* -------------------------------------------------------------------------- */

// Parses a 'Platform' struct from a platform identifier. There are potentially
// multiple valid identifiers for any given platform due to schema differences
// across Godot versions.
func Parse(input string) (Platform, error) { //nolint:cyclop
	if input == "" {
		return Platform{}, ErrMissingPlatform
	}

	switch strings.ToLower(strings.TrimSpace(input)) {
	// Linux
	case "x11.32", "linux.x86_32":
		return Platform{I386, Linux}, nil
	case "x11.64", "linux.x86_64":
		return Platform{Amd64, Linux}, nil

	// Linux (mono builds)
	case "linux_x86_32":
		return Platform{I386, Linux}, nil
	case "linux_x86_64":
		return Platform{Amd64, Linux}, nil

	// MacOS - Note that the supported architectures between 'osx.fat' and
	// 'osx.universal' are *NOT* the same. It's important to maintain the
	// 'Version' alongside this result so that the architectures can be
	// correctly determined.
	case "osx.64":
		return Platform{Amd64, MacOS}, nil
	case "macos.universal", "osx.fat", "osx.universal":
		return Platform{Universal, MacOS}, nil

	// Windows
	case "win32":
		return Platform{I386, Windows}, nil
	case "win64":
		return Platform{Amd64, Windows}, nil

	default:
		return Platform{}, fmt.Errorf("%w: '%s'", ErrUnrecognizedPlatform, input)
	}
}

/* --------------------------- Function: MustParse -------------------------- */

// Parses an input string as a 'Platform' specification but panics if it would
// fail.
func MustParse(input string) Platform {
	platform, err := Parse(input)
	if err != nil {
		panic(err)
	}

	return platform
}

/* -------------------------------------------------------------------------- */
/*                              Function: Format                              */
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
func Format(p Platform, v godot.Version) (string, error) {
	// Use the 'Platform' validation in 'New' prior to formatting; a default
	// 'Platform' is not valid, which is why this check is required.
	if _, err := New(p.OS, p.Arch); err != nil {
		return "", err
	}

	switch p.OS {
	case Linux:
		return formatLinux(p.Arch, v)
	case MacOS:
		return formatMacOS(p.Arch, v)
	case Windows:
		return formatWindows(p.Arch, v)

	case 0:
		return "", ErrMissingOS
	default:
		return "", ErrUnrecognizedOS
	}
}

/* -------------------------- Function: formatLinux ------------------------- */

// Given an architecture, returns the Linux platform identifier used by Godot
// executable names.
func formatLinux(a Arch, v godot.Version) (string, error) { //nolint:cyclop
	if a == 0 {
		return "", ErrMissingArch
	}

	var p string

	switch {
	// v1-v2 not supported
	case v.Major() < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", godot.ErrUnsupportedVersion, v)
	// v3
	case v.Major() < 4: //nolint:gomnd
		// 'linux_headless.64' and 'linux_server.64' flavors introduced in v3.1
		// are not supported.
		switch a {
		case I386:
			p = "x11.32"
		case Amd64:
			p = "x11.64"

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// v4.0-dev.20210727 - Godot v4.0-alpha14
	case v.CompareNormal(godot.V4()) == 0 && reV4LinuxLabelsWithoutX86.MatchString(v.Label()):
		switch a {
		case I386:
			p = "linux.32"
		case Amd64:
			p = "linux.64"

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// v4.0-alpha15+
	default:
		switch a {
		case I386:
			p = "linux.x86_32"
		case Amd64:
			p = "linux.x86_64"

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}

	// All "mono"-flavored builds have the '.' rune replaced by a '_' rune.
	if v.IsMono() {
		p = strings.ReplaceAll(p, ".", "_")
	}

	return p, nil
}

/* -------------------------- Function: formatMacOS ------------------------- */

// Given an architecture, returns the macOS platform identifier used by Godot
// executable names.
//
// NOTE: This is rather convoluted; consider a better way of organizing this
// logic.
func formatMacOS(a Arch, v godot.Version) (string, error) { //nolint:cyclop
	if a == 0 {
		return "", ErrMissingArch
	}

	switch {
	// v1 - v2 not supported
	case v.Major() < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", godot.ErrUnsupportedVersion, v)
	// v3.0 - v3.0.6
	case v.Major() == 3 && v.Minor() < 1:
		switch a {
		case I386, Amd64:
			return "osx.fat", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// v3.1 - v3.2.4-beta2
	// NOTE: Because v3.2.4 labels are only "beta" and "rc" *and* "beta"
	// versions do not exceed 6, lexicographic  sorting works.
	case v.Major() == 3 && v.Minor() <= 2 && (v.Patch() < 4 || v.Patch() == 4 && v.Label() <= "beta2"):
		switch a {
		case Amd64:
			return "osx.64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// v3.2.4-beta3 - v4.0-alpha12
	case v.CompareNormal(godot.V4()) < 0 ||
		(v.CompareNormal(godot.V4()) == 0 && reV4MacOSLabelsWithOSXUniversal.MatchString(v.Label())):
		switch a {
		case Amd64, Arm64:
			return "osx.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	// v4.0-alpha13+
	default:
		switch a {
		case Amd64, Arm64:
			return "macos.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}
}

/* ------------------------- Function: formatWindows ------------------------ */

// Given an architecture, returns the Windows platform identifier used by Godot
// executable names.
func formatWindows(a Arch, v godot.Version) (string, error) {
	if a == 0 {
		return "", ErrMissingArch
	}

	switch {
	// v1-v2 not supported
	case v.Major() < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: '%s'", godot.ErrUnsupportedVersion, v)
	// v3+
	default:
		switch a {
		case I386:
			return "win32", nil
		case Amd64:
			return "win64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArch, a)
		}
	}
}
