package commands

import (
	"log"
	"os"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print installed versions of Godot.
func NewLs() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Category: "Utilities",

		Usage:     "print the path and version of all of the installed versions of Godot",
		UsageText: "gdenv ls",

		Action: func(_ *cli.Context) error {
			// Ensure 'Store' layout
			s, err := store.Path()
			if err != nil {
				return cli.Exit(err, 1)
			}

			if err := store.Init(s); err != nil {
				return cli.Exit(err, 1)
			}

			results, err := ls(s)
			if err != nil {
				return cli.Exit(err, 1)
			}

			if wd, err := os.Getwd(); err == nil {
				if version, err := pin.Read(wd); err == nil {
					log.Printf("🤖 Currently active version: %s\n\n", version)
				}
			}

			log.Printf("Installed versions:\n\n%s\n", strings.Join(results, "\n"))

			return nil
		},
	}
}

func ls(s string) ([]string, error) {
	versions, err := store.List(s)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(versions))
	for i, v := range versions {
		out[i] = v.String()
	}

	return out, nil
}
