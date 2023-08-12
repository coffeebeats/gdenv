package main

import (
	"errors"
	"log"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

var (
	ErrGodotNotFound = errors.New("godot not found")
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
			// Determine 'path' option
			path, err := resolvePath(c)
			if err != nil {
				return err
			}

			// Ensure 'Store' layout
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			if err := store.Init(storePath); err != nil {
				return err
			}

			toolPath, err := which(storePath, path)
			if err != nil {
				return err
			}

			log.Println(toolPath)

			return nil
		},
	}
}

func which(storePath, pinPath string) (string, error) {
	path, err := pin.Resolve(pinPath)
	if err != nil {
		if !errors.Is(err, pin.ErrFileNotFound) {
			return "", err
		}
	}

	// No pin file was found yet, so check globally.
	if path == "" {
		path = storePath
	}

	version, err := pin.Read(path)
	if err != nil || !store.Has(storePath, version) {
		return "", ErrGodotNotFound
	}

	return store.ToolPath(storePath, version)
}
