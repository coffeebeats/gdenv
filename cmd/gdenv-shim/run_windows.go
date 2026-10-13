//go:build windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func run(binary string, args ...string) error {
	cmd := exec.Command(binary, args...) //nolint:gosec,noctx

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //nolint:exhaustruct_v5

	return cmd.Run()
}
