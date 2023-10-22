package main

import (
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

/* ----------------------------- Function: NewLs ---------------------------- */

// A 'urfave/cli' command to print installed versions of Godot.
func NewLs() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Category: "Utilities",

		Usage:     "print the path and version of all of the installed versions of Godot",
		UsageText: "gdenv ls",

		Action: func(c *cli.Context) error {
			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			executables, err := store.Executables(c.Context, storePath)
			if err != nil {
				return err
			}

			if len(executables) == 0 {
				return nil
			}

			printedActive, err := printActiveVersion(c.Context, storePath)
			if err != nil {
				return err
			}

			printedGlobal, err := printGlobalVersion(c.Context, storePath)
			if err != nil {
				return err
			}

			if printedActive || printedGlobal {
				log.Print("")
			}

			log.Printf("Installed versions (%s):\n", storePath)

			for _, ex := range executables {
				platformLabel, err := platform.Format(ex.Artifact.Platform(), ex.Artifact.Version())
				if err != nil {
					return err
				}

				log.Printf("  %s (%s)\n", ex.Artifact.Version(), platformLabel)
			}

			return nil
		},
	}
}

/* ---------------------- Function: printActiveVersion ---------------------- */

func printActiveVersion(ctx context.Context, storePath string) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	v, err := pin.VersionAt(ctx, storePath, wd)
	if err != nil {
		if !errors.Is(err, pin.ErrMissingPin) {
			return false, err
		}

		return false, nil
	}

	// Define the host 'Platform'.
	p, err := platform.Detect()
	if err != nil {
		return false, err
	}

	platformLabel, err := platform.Format(p, v)
	if err != nil {
		return false, err
	}

	log.Printf("🤖 Currently active version: %s (%s)", v, platformLabel)

	return true, nil
}

/* ---------------------- Function: printGlobalVersion ---------------------- */

func printGlobalVersion(ctx context.Context, storePath string) (bool, error) {
	v, err := pin.VersionAt(ctx, storePath, storePath)
	if err != nil {
		if !errors.Is(err, pin.ErrMissingPin) {
			return false, err
		}

		return false, nil
	}

	log.Printf("🌎 System default version: %s", v)

	return true, nil
}
