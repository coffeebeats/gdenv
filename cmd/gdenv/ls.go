package main

import (
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
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

		Aliases: []string{"list"},

		Usage:     "print the path and version of all of the installed versions of Godot",
		UsageText: "gdenv ls",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "list executable and source code versions",
			},
			&cli.BoolFlag{
				Name:    "source",
				Aliases: []string{"s", "src"},
				Usage:   "list source code versions",
			},
		},

		Action: func(c *cli.Context) error {
			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			src, all := c.Bool("source"), c.Bool("all")

			if !src {
				ok, err := printGlobalVersion(c.Context, storePath)
				if err != nil {
					return err
				}

				if ok {
					log.Print("")
				}
			}

			if src || all {
				if err := printSources(c.Context, storePath); err != nil {
					return err
				}
			}

			if !src || all {
				if err := printExecutables(c.Context, storePath); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

/* ---------------------- Function: printGlobalVersion ---------------------- */

func printGlobalVersion(ctx context.Context, storePath string) (bool, error) {
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

	log.Printf("🌎 System default version: %s (%s)", v, platformLabel)

	ok, err := store.Has(storePath, executable.New(v, p))
	if err != nil {
		return false, err
	}

	if !ok {
		log.Warn("system default version not installed")
	}

	return true, nil
}

/* ----------------------- Function: PrintExecutables ----------------------- */

func printExecutables(ctx context.Context, storePath string) error {
	executables, err := store.Executables(ctx, storePath)
	if err != nil {
		return err
	}

	if len(executables) == 0 {
		return nil
	}

	log.Printf("Installed executable versions (%s):", storePath)

	for _, ex := range executables {
		platformLabel, err := platform.Format(ex.Artifact.Platform(), ex.Artifact.Version())
		if err != nil {
			return err
		}

		log.Printf("  %s (%s)", ex.Artifact.Version(), platformLabel)
	}

	return nil
}

/* ------------------------- Function: PrintSources ------------------------- */

func printSources(ctx context.Context, storePath string) error {
	sources, err := store.Sources(ctx, storePath)
	if err != nil {
		return err
	}

	if len(sources) == 0 {
		return nil
	}

	log.Printf("Installed source code versions (%s):", storePath)

	for _, src := range sources {
		log.Printf("  %s", src.Artifact.Inner.Version())
	}

	return nil
}
