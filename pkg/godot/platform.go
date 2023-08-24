package godot

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
)

var (
	ErrMissingArchInput      = errors.New("missing architecture input")
	ErrMissingOSInput        = errors.New("missing OS input")
	ErrUnrecognizedArchInput = errors.New("unrecognized architecture input")
	ErrUnrecognizedOSInput   = errors.New("unrecognized OS input")
	ErrUnsupportedArchInput  = errors.New("unsupported architecture input")
	ErrUnsupportedOSInput    = errors.New("unsupported OS input")

	versionGodot4          = Version{major: 4, minor: 0} //nolint:gomnd
	versionMacOSArmSupport = Version{3, 2, 4, "beta3"}

	// NOTE: Unfortunately, there isn't a clean logic to how these versions are
	// labeled. Rather than implementing some rules based on label parsing, just
	// maintain a list of labels still using the 'osx.universal' identifier.
	versionsGodot4WithOSXUniversal = []string{
		"alpha1",
		"alpha2",
		"alpha3",
		"alpha4",
		"alpha5",
		"alpha6",
		"alpha7",
		"alpha8",
		"alpha9",
		"alpha10",
		"alpha11",
		"alpha12",

		"dev.20210727",
		"dev.20210811",
		"dev.20210820",
		"dev.20210916",
		"dev.20210924",
		"dev.20211004",
		"dev.20211015",
		"dev.20211027",
		"dev.20211108",
		"dev.20211117",
		"dev.20211210",
		"dev.20220105",
		"dev.20220118",
	}
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
		return 0, ErrMissingOSInput

	case "darwin", "macos", "osx":
		return MacOS, nil

	case "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
		return Linux, nil

	case "win", "windows":
		return Windows, nil

	default:
		return 0, fmt.Errorf("%w: %s", ErrUnrecognizedOSInput, input)
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
		return 0, ErrMissingArchInput

	case "386", "i386", "x86", "x86_32":
		return I386, nil

	case "amd64", "x86_64", "x86-64":
		return Amd64, nil

	case "arm64", "arm64be":
		return Arm64, nil

	case "fat", "universal":
		return Universal, nil

	default:
		return 0, fmt.Errorf("%w: %s", ErrUnrecognizedArchInput, input)
	}
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Platform                              */
/* -------------------------------------------------------------------------- */

// A platform specification representing a target to run the Godot editor on.
type Platform struct {
	arch Arch
	os   OS
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

	platform.arch, platform.os = arch, oS

	return platform, nil
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
	switch p.os {
	case Linux:
		return formatLinuxPlatform(p.arch, v)
	case MacOS:
		return formatMacOSPlatform(p.arch, v)
	case Windows:
		return formatWindowsPlatform(p.arch, v)

	case 0:
		return "", ErrMissingOSInput
	default:
		return "", ErrUnrecognizedOSInput
	}
}

/* ---------------------- Function: formatLinuxPlatform --------------------- */

func formatLinuxPlatform(a Arch, v Version) (string, error) {
	if a == 0 {
		return "", ErrMissingArchInput
	}

	switch {
	// Godot v1-v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: %s", ErrUnsupportedVersion, v)
	// Godot v3
	case v.major < 4: //nolint:gomnd
		// 'linux_headless.64' and 'linux_server.64' flavors introduced in v3.1
		// are not supported.
		switch a {
		case I386:
			return "x11.32", nil
		case Amd64:
			return "x11.64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	// Godot v4+
	default:
		switch a {
		case I386:
			return "linux.x86_32", nil
		case Amd64:
			return "linux.x86_64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	}
}

/* ---------------------- Function: formatMacOSPlatform --------------------- */

func formatMacOSPlatform(a Arch, v Version) (string, error) {
	if a == 0 {
		return "", ErrMissingArchInput
	}

	switch {
	// Godot v1 - v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: %s", ErrUnsupportedVersion, v)
	// Godot v3.0 - v3.0.6
	case v.major == 3 && v.minor < 1:
		switch a {
		case I386, Amd64:
			return "osx.fat", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	// Godot v3.1 - v3.2.4-beta2
	// NOTE: Because v3.2.4 labels are only "beta" and "rc" *and* "beta"
	// versions do not exceed 6, lexicographic  sorting works.
	case v.major == 3 && v.minor <= 2 && (v.patch < 4 || v.patch == 4 && v.label <= "beta2"):
		switch a {
		case Amd64:
			return "osx.64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	// Godot v3.2.4-beta3 - v4.0-alpha12
	case v.Compare(versionGodot4) < 0 ||
		// An O(n) check here is fine - the data is small and this will run at
		// most once per CLI invocation.
		(v.Compare(versionGodot4) == 0 && slices.Contains(versionsGodot4WithOSXUniversal, v.label)):
		switch a {
		case Amd64, Arm64:
			return "osx.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	// Godot v4.0-alpha13+
	default:
		switch a {
		case Amd64, Arm64:
			return "macos.universal", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	}
}

/* --------------------- Function: formatWindowsPlatform -------------------- */

func formatWindowsPlatform(a Arch, v Version) (string, error) {
	if a == 0 {
		return "", ErrMissingArchInput
	}

	switch {
	// Godot v1-v2 not supported
	case v.major < 3: //nolint:gomnd
		return "", fmt.Errorf("%w: %s", ErrUnsupportedVersion, v)
	// Godot v3+
	default:
		switch a {
		case I386:
			return "win32", nil
		case Amd64:
			return "win64", nil

		default:
			return "", fmt.Errorf("%w: %v", ErrUnsupportedArchInput, a)
		}
	}
}
