package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/store"
)

// A 'urfave/cli' command to delete a cached version of Godot.
func NewUninstall() *cli.Command {
	return &cli.Command{
		Name:     "uninstall",
		Category: "Install",

		Usage:     "Remove the specified version of Godot from the 'gdenv' download cache",
		UsageText: "gdenv uninstall [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "uninstall all versions of Godot (ignores source code without '-s')",
			},
			&cli.BoolFlag{
				Name:    "source",
				Aliases: []string{"s", "src"},
				Usage:   "uninstall source code versions",
			},
		},

		Action: func(c *cli.Context) error {
			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			src, all := c.Bool("source"), c.Bool("all")

			// Uninstall all versions.
			switch {
			case src && all:
				return uninstallAllSources(c.Context, storePath)
			case !src && all:
				return uninstallAllExecutables(c.Context, storePath)
			}

			// Uninstall a specific version.

			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil && !c.Bool("all") {
				return UsageError{ctx: c, err: err}
			}

			log.Infof("uninstalling version: %s", v)

			switch {
			case src:
				return store.Remove(storePath, source.New(v))
			default:
				return uninstallExecutable(storePath, v)
			}
		},
	}
}

/* -------------------- Function: uninstallAllExecutables ------------------- */

func uninstallAllExecutables(ctx context.Context, storePath string) error {
	ee, err := store.Executables(ctx, storePath)
	if err != nil {
		return err
	}

	if len(ee) == 0 {
		return nil
	}

	log.Info("removing all installed executable versions")

	return store.Clear(storePath)
}

/* ---------------------- Function: uninstallAllSources --------------------- */

func uninstallAllSources(ctx context.Context, storePath string) error {
	ss, err := store.Sources(ctx, storePath)
	if err != nil {
		return err
	}

	if len(ss) == 0 {
		return nil
	}

	log.Info("removing all installed source code versions")

	return store.Clear(storePath)
}

/* ---------------------- Function: uninstallExecutable --------------------- */

func uninstallExecutable(storePath string, v version.Version) error {
	// Define the host 'Platform'.
	p, err := platform.Detect()
	if err != nil {
		return err
	}

	// Define the target 'Executable'.
	ex := executable.New(v, p)

	return store.Remove(storePath, ex)
}
