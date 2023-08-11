package commands

import (
	"errors"
	"log"
	"os"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print the path to the effective Godot binary.
func NewWhich() *cli.Command {
	return &cli.Command{
		Name:     "which",
		Category: "Utilities",

		Usage:     "print the path to the Godot executable which would be used in the specified directory",
		UsageText: "gdenv which [OPTIONS]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Check the specified `PATH`",
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

			// Determine 'path' option
			var path string
			switch {
			case c.Bool("path"):
				path = c.String("path")
			default:
				p, err := os.Getwd()
				if err != nil {
					return cli.Exit(err, 1)
				}

				path = p
			}

			f, err := which(s, path)
			if err != nil {
				return cli.Exit(err, 1)
			}

			log.Println(f)

			return nil
		},
	}
}

func which(s, path string) (string, error) {
	path, err := pin.Resolve(path)
	if err != nil {
		if !errors.Is(err, pin.ErrFileNotFound) {
			return "", err
		}
	}

	// No pin file was found yet
	if path == "" {
		// Check globally
		p, err := store.Path()
		if err != nil {
			return "", nil
		}

		path = p
	}

	version, err := pin.Read(path)
	if err != nil {
		return "", errors.New("godot not found")
	}

	return store.Find(s, version)
}
