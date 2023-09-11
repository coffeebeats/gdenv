package platform

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingArch      = errors.New("missing architecture")
	ErrUnrecognizedArch = errors.New("unrecognized architecture")
	ErrUnsupportedArch  = errors.New("unsupported architecture")
)

/* -------------------------------------------------------------------------- */
/*                                 Enum: Arch                                 */
/* -------------------------------------------------------------------------- */

// CPU architectures which the Godot project provides prebuilt editor binaries
// for.
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
		return 0, ErrMissingArch

	case "386", "i386", "x86", "x86_32":
		return I386, nil

	case "amd64", "x86_64", "x86-64":
		return Amd64, nil

	case "arm64", "arm64be":
		return Arm64, nil

	case "fat", "universal":
		return Universal, nil

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
