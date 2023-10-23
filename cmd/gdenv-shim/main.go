package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
)

var (
	ErrMissingPin   = errors.New("no version selected; try setting a version pin with 'gdenv pin'")
	ErrNotInstalled = errors.New("pinned version not installed; try installing with 'gdenv install'")
)

func main() {
	var exitCode int
	defer func() {
		if err := recover(); err != nil {
			exitCode = 1

			log.Println(err)
		}

		os.Exit(exitCode)
	}()

	// Don't report timestamp in logs.
	log.SetFlags(0)

	if err := execute(context.Background()); err != nil {
		panic(err)
	}
}

/* ---------------------------- Function: execute --------------------------- */

// execute replaces the current process with the cached version of Godot
// specified by a local or global pin.
func execute(ctx context.Context) error {
	p, err := platform.Detect()
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	storePath, err := store.Path()
	if err != nil {
		return err
	}

	binary, err := install.Which(ctx, storePath, p, wd)
	if err != nil {
		if errors.Is(err, pin.ErrMissingPin) {
			return ErrMissingPin
		}

		if errors.Is(err, install.ErrNotInstalled) {
			return ErrNotInstalled
		}

		return err
	}

	switch p.OS {
	case platform.Windows:
		cmd := exec.Command(binary, os.Args[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	default:
		return syscall.Exec(
			binary,
			os.Args[1:],
			os.Environ(),
		)
	}
}
