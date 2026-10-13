//go:build !windows

package main

import (
	"os"
	"syscall"
)

/* ------------------------------ Function: run ----------------------------- */

// run replaces the shim's process with Godot's, so it returns only if that
// fails.
func run(binary string, args ...string) (int, error) {
	return 0, syscall.Exec( //nolint:gosec
		binary,
		append([]string{binary}, args...),
		os.Environ(),
	)
}
