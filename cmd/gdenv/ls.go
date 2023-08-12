package main

import (
	"log"
	"os"
	"strings"

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

		Action: func(_ *cli.Context) error {
			// Ensure 'Store' layout
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			if err := store.Init(storePath); err != nil {
				return err
			}

			results, err := ls(storePath)
			if err != nil {
				return err
			}

			if wd, err := os.Getwd(); err == nil {
				if version, err := pin.Read(wd); err == nil {
					log.Printf("🤖 Currently active version: %s\n\n", version.Canonical())
				}
			}

			log.Printf("Installed versions:\n\n%s\n", strings.Join(results, "\n"))

			return nil
		},
	}
}

/* ------------------------------ Function: ls ------------------------------ */

func ls(path string) ([]string, error) {
	versions, err := store.Versions(path)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(versions))
	for i, v := range versions {
		out[i] = v.String()
	}

	return out, nil
}
