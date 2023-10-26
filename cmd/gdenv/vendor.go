package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to download and cache a specific version of the Godot
// source code.
func NewVendor() *cli.Command {
	return &cli.Command{
		Name:     "vendor",
		Category: "Install",

		Usage:     "download and cache a specific version of Godot source code",
		UsageText: "gdenv install [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry",
			},
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "download the source code into `OUT` (overwrites conflicting files; defaults to './godot')",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "determine the version from the pinned `PATH` (ignores the global pin)",
			},
		},

		Action: func(c *cli.Context) error {
			v, err := resolveVersionFromArgOrPath(c)
			if err != nil {
				return err
			}

			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			if err := installSource(c.Context, storePath, v, c.Bool("force")); err != nil {
				return err
			}

			if !c.IsSet("path") {
				return nil
			}

			out := c.String("out")
			if out == "" {
				out = "./godot"
			}

			return vendor(c.Context, v, storePath, out)
		},
	}
}

/* ------------------------- Function: installSource ------------------------ */

// Installs the specified version of the source code  to the store, but only if
// needed.
func installSource(ctx context.Context, storePath string, v version.Version, force bool) error {
	log.Infof("installing version: %s", v)

	// Ensure the store exists.
	if err := store.Touch(storePath); err != nil {
		return err
	}

	// Define the target 'Source'.
	src := source.New(v)

	ok, err := store.Has(storePath, src)
	if err != nil {
		return err
	}

	if ok && !force {
		log.Info("skipping installation; version already found")

		return nil
	}

	if err := install.Source(ctx, storePath, src); err != nil {
		return err
	}

	log.Infof("successfully installed version: %s", src.Version())

	return nil
}

/* ---------------------------- Function: vendor ---------------------------- */

// Extracts the cached source code folder into the specified 'out' path.
func vendor(ctx context.Context, v version.Version, storePath, out string) error {
	src := source.Archive{Artifact: source.New(v)}

	srcPath, err := store.Source(storePath, src.Artifact)
	if err != nil {
		return err
	}

	localSrcArchive := artifact.Local[source.Archive]{Artifact: src, Path: srcPath}
	if err := archive.Extract(ctx, localSrcArchive, out); err != nil {
		return err
	}

	log.Infof("successfully vendored version %s: %s", v, srcPath)

	return nil
}
