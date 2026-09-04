package update

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	nameBinary = "gdenv"
	nameShim   = "godot"

	extensionExe   = ".exe"
	extensionTarGz = ".tar.gz"
	extensionZip   = ".zip"

	osLinux   = "linux"
	osMacOS   = "macos"
	osWindows = "windows"

	archARM64 = "arm64"
	archX8664 = "x86_64"
)

var ErrUnsupportedTarget = errors.New("unsupported target")

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

// A Target identifies which published release archive applies to a host, using
// the naming scheme defined by '.goreleaser.yaml'.
//
// NOTE: These names are *not* the same as those in 'pkg/godot/platform', which
// describe Godot's own release artifacts.
type Target struct {
	OS   string
	Arch string
}

/* ------------------------- Function: DetectTarget ------------------------- */

// DetectTarget determines the release target for the running binary.
//
// NOTE: This deliberately ignores the 'GDENV_OS' and 'GDENV_ARCH' environment
// variables. Those select which *Godot* build to install; the 'gdenv' binary
// already knows what it was compiled for.
func DetectTarget() (Target, error) {
	return newTarget(runtime.GOOS, runtime.GOARCH)
}

/* --------------------------- Function: newTarget -------------------------- */

// newTarget implements 'DetectTarget' for an explicit build target, so that
// every published combination can be exercised from a single host.
func newTarget(goos, goarch string) (Target, error) {
	var target Target

	switch goos {
	case "darwin":
		target.OS = osMacOS
	case "linux":
		target.OS = osLinux
	case "windows":
		target.OS = osWindows
	default:
		return target, fmt.Errorf("%w: no prebuilt binaries for operating system: %s", ErrUnsupportedTarget, goos)
	}

	switch goarch {
	case "amd64":
		target.Arch = archX8664
	case "arm64":
		target.Arch = archARM64
	default:
		return target, fmt.Errorf("%w: no prebuilt binaries for CPU architecture: %s", ErrUnsupportedTarget, goarch)
	}

	// NOTE: Only the combinations built by '.goreleaser.yaml' are published;
	// Windows on ARM64 is not among them.
	if target.OS == osWindows && target.Arch == archARM64 {
		return Target{}, fmt.Errorf(
			"%w: no prebuilt '%s' binaries for operating system: %s",
			ErrUnsupportedTarget, target.Arch, target.OS,
		)
	}

	return target, nil
}

/* -------------------------- Method: ArchiveName --------------------------- */

// ArchiveName returns the filename of the release archive for this target at
// the specified version.
//
// NOTE: This must match the 'archives.name_template' value defined in
// '.goreleaser.yaml'; 'TestTargetsMatchGoreleaser' guards that.
func (t Target) ArchiveName(v string) string {
	return fmt.Sprintf("%s-%s-%s-%s%s", nameBinary, v, t.OS, t.Arch, t.archiveExtension())
}

/* ---------------------------- Method: Binaries ---------------------------- */

// Binaries returns the names of the executables contained within a release
// archive.
//
// NOTE: The order is significant. The 'gdenv' binary is last because it is the
// one currently running, making it the most likely to fail to be replaced; see
// 'Apply' for how a partial replacement is rolled back.
func (t Target) Binaries() []string {
	if t.OS == osWindows {
		return []string{nameShim + extensionExe, nameBinary + extensionExe}
	}

	return []string{nameShim, nameBinary}
}

/* ------------------------ Method: archiveExtension ------------------------ */

func (t Target) archiveExtension() string {
	if t.OS == osWindows {
		return extensionZip
	}

	return extensionTarGz
}
