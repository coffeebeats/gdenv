package platform

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingOS      = errors.New("missing OS")
	ErrUnrecognizedOS = errors.New("unrecognized OS")
)

/* -------------------------------------------------------------------------- */
/*                                  Enum: OS                                  */
/* -------------------------------------------------------------------------- */

// Operating systems which the Godot project provides prebuilt editor binaries
// for.
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
		return 0, ErrMissingOS

	case "darwin", "macos", "osx":
		return MacOS, nil

	case "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
		return Linux, nil

	case "win", "windows":
		return Windows, nil

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
