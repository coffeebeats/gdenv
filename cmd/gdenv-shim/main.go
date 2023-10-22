package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"
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
		return err
	}

	return syscall.Exec(
		binary,
		append([]string{binary}, os.Args[1:]...),
		os.Environ(),
	)
}
