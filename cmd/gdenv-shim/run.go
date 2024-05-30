//go:build !windows

package main

import (
	"os"
	"syscall"
)

func run(binary string, args ...string) error {
	return syscall.Exec( //nolint:gosec
		binary,
		append([]string{binary}, args...),
		os.Environ(),
	)
}
