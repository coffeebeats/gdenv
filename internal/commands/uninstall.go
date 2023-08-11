package commands

import (
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to delete a cached version of Godot.
func NewUninstall() *cli.Command {
	return &cli.Command{
		Name:     "uninstall",
		Category: "Install",

		Usage:     "remove the specified version of Godot from the gdenv download cache",
		UsageText: "gdenv uninstall [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "uninstall all versions of Godot in the cache",
			},
		},

		Action: func(c *cli.Context) error {
			// Ensure 'Store' layout
			s, err := store.Path()
			if err != nil {
				return cli.Exit(err, 1)
			}

			if err := store.Init(s); err != nil {
				return cli.Exit(err, 1)
			}

			if c.Bool("all") {
				versions, err := store.List(s)
				if err != nil {
					return cli.Exit(err, 1)
				}

				for _, v := range versions {
					if err := uninstall(s, v); err != nil {
						return cli.Exit(err, 1)
					}
				}
			} else {
				// Validate arguments
				version, err := godot.ParseVersion(c.Args().First())
				if err != nil && !c.Bool("all") {
					return cli.Exit(err, 1)
				}

				if err := uninstall(s, version); err != nil {
					return cli.Exit(err, 1)
				}
			}

			return nil
		},
	}
}

func uninstall(s string, v godot.Version) error {
	p, err := store.Find(s, v)
	if err != nil {
		return err
	}

	// Remove the tool and the parent directory
	return os.RemoveAll(filepath.Base(p))
}
