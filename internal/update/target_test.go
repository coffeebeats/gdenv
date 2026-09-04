package update

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                            Function: TestNewTarget                         */
/* -------------------------------------------------------------------------- */

func TestNewTarget(t *testing.T) {
	tests := []struct {
		goos, goarch string
		want         Target
		err          error
	}{
		// Published targets.
		{goos: "darwin", goarch: "amd64", want: Target{OS: osMacOS, Arch: archX8664}},
		{goos: "darwin", goarch: "arm64", want: Target{OS: osMacOS, Arch: archARM64}},
		{goos: "linux", goarch: "amd64", want: Target{OS: osLinux, Arch: archX8664}},
		{goos: "linux", goarch: "arm64", want: Target{OS: osLinux, Arch: archARM64}},
		{goos: "windows", goarch: "amd64", want: Target{OS: osWindows, Arch: archX8664}},

		// Unpublished combination.
		{goos: "windows", goarch: "arm64", err: ErrUnsupportedTarget},

		// Unsupported inputs.
		{goos: "freebsd", goarch: "amd64", err: ErrUnsupportedTarget},
		{goos: "linux", goarch: "386", err: ErrUnsupportedTarget},
		{goos: "linux", goarch: "riscv64", err: ErrUnsupportedTarget},
	}

	for _, tc := range tests {
		t.Run(tc.goos+"_"+tc.goarch, func(t *testing.T) {
			got, err := newTarget(tc.goos, tc.goarch)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			if tc.err != nil {
				return
			}

			if got != tc.want {
				t.Fatalf("target: got %+v, want %+v", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                          Function: TestArchiveName                         */
/* -------------------------------------------------------------------------- */

func TestArchiveName(t *testing.T) {
	tests := []struct {
		target Target
		want   string
	}{
		{Target{OS: osMacOS, Arch: archX8664}, "gdenv-v1.2.3-macos-x86_64.tar.gz"},
		{Target{OS: osMacOS, Arch: archARM64}, "gdenv-v1.2.3-macos-arm64.tar.gz"},
		{Target{OS: osLinux, Arch: archX8664}, "gdenv-v1.2.3-linux-x86_64.tar.gz"},
		{Target{OS: osLinux, Arch: archARM64}, "gdenv-v1.2.3-linux-arm64.tar.gz"},
		{Target{OS: osWindows, Arch: archX8664}, "gdenv-v1.2.3-windows-x86_64.zip"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.target.ArchiveName("v1.2.3"); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                      Function: TestArchiveNameVersions                     */
/* -------------------------------------------------------------------------- */

// TestArchiveNameVersions pins how the version is embedded in an asset name.
// The version is interpolated verbatim, matching goreleaser's
// '{{ .ProjectName }}-v{{ .Version }}-...' template, where '.Version' is the
// git tag with its leading 'v' removed.
func TestArchiveNameVersions(t *testing.T) {
	target := Target{OS: osLinux, Arch: archX8664}

	tests := []struct{ version, want string }{
		// The version currently published.
		{version: "v0.6.35", want: "gdenv-v0.6.35-linux-x86_64.tar.gz"},

		// A zero patch component is preserved rather than truncated; note that
		// 'pkg/godot/version' would render this as 'v1.2', which is why this
		// package does not reuse it.
		{version: "v1.2.0", want: "gdenv-v1.2.0-linux-x86_64.tar.gz"},
		{version: "v0.0.0", want: "gdenv-v0.0.0-linux-x86_64.tar.gz"},

		// Multi-digit components must not be padded or reordered.
		{version: "v10.20.30", want: "gdenv-v10.20.30-linux-x86_64.tar.gz"},

		// A pre-release suffix contains '-', the same character which separates
		// the name's fields. Asset names are only ever constructed, never parsed
		// back apart, so this is unambiguous in the direction that matters.
		{version: "v1.2.3-rc.1", want: "gdenv-v1.2.3-rc.1-linux-x86_64.tar.gz"},
	}

	for _, tc := range tests {
		t.Run(tc.version, func(t *testing.T) {
			if got := target.ArchiveName(tc.version); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: TestBinaries                           */
/* -------------------------------------------------------------------------- */

func TestBinaries(t *testing.T) {
	// The running 'gdenv' binary must be replaced last; see 'Apply'.
	if got := (Target{OS: osLinux, Arch: archX8664}).Binaries(); got[len(got)-1] != nameBinary {
		t.Fatalf("expected '%s' to be replaced last, got %v", nameBinary, got)
	}

	if got := (Target{OS: osWindows, Arch: archX8664}).Binaries(); got[len(got)-1] != nameBinary+extensionExe {
		t.Fatalf("expected '%s' to be replaced last, got %v", nameBinary+extensionExe, got)
	}
}

/* -------------------------------------------------------------------------- */
/*                    Function: TestTargetsMatchGoreleaser                    */
/* -------------------------------------------------------------------------- */

// TestTargetsMatchGoreleaser guards the duplication between this package and
// '.goreleaser.yaml': the release asset naming lives in Go here and in the
// install scripts, so a change to what is published must not silently diverge
// from what this package will download.
func TestTargetsMatchGoreleaser(t *testing.T) {
	// Given: The release configuration which decides what is published.
	b, err := os.ReadFile("../../.goreleaser.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// NOTE: Line endings are normalized; this repository checks out CRLF on
	// Windows, which would otherwise defeat the '$' anchors below.
	config := strings.ReplaceAll(string(b), "\r\n", "\n")

	// When: Every 'goos_goarch' target listed in that config is enumerated.
	matches := regexp.MustCompile(`(?m)^\s+- ([a-z0-9]+)_([a-z0-9]+)$`).FindAllStringSubmatch(config, -1)
	if len(matches) == 0 {
		t.Fatal("found no build targets in '.goreleaser.yaml'")
	}

	published := make(map[string]struct{})

	for _, m := range matches {
		published[m[1]+"_"+m[2]] = struct{}{}
	}

	// Then: Every published target is one this package can resolve.
	for target := range published {
		goos, goarch, _ := strings.Cut(target, "_")

		if _, err := newTarget(goos, goarch); err != nil {
			t.Errorf("'%s' is published but not supported by this package: %v", target, err)
		}
	}

	// Then: Every target this package claims to support is published.
	for _, goos := range []string{"darwin", "linux", "windows"} {
		for _, goarch := range []string{"amd64", "arm64"} {
			_, err := newTarget(goos, goarch)

			_, isPublished := published[goos+"_"+goarch]

			if err == nil && !isPublished {
				t.Errorf("'%s_%s' is supported by this package but is not published", goos, goarch)
			}

			if err != nil && isPublished {
				t.Errorf("'%s_%s' is published but rejected by this package", goos, goarch)
			}
		}
	}

	// Then: The archive naming template still produces the names built above.
	for _, fragment := range []string{
		"{{ .ProjectName }}-v{{ .Version }}-",
		`{{- if eq .Os "darwin" }}macos`,
		`{{- if eq .Arch "amd64" }}x86_64`,
		"name_template: checksums.txt",
	} {
		if !strings.Contains(config, fragment) {
			t.Errorf(
				"'.goreleaser.yaml' no longer contains %q; "+
					"the naming in 'target.go' may need updating",
				fragment,
			)
		}
	}
}
