package main

import (
	"context"
	"errors"
	"io/fs"
	"log"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
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
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			// Define the host 'Platform'.
			platform, err := platform.Detect()
			if err != nil {
				return err
			}

			toolPath, err := which(c.Context, storePath, path, platform)
			if err != nil {
				return err
			}

			log.Println(toolPath)

			return nil
		},
	}
}

func which(ctx context.Context, storePath, pinPath string, p platform.Platform) (string, error) {
	path, err := pin.Resolve(ctx, pinPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
	}

	// No pin file was found yet, so check globally.
	if path == "" {
		path = storePath
	}

	v, err := pin.Read(path)
	if err != nil {
		return "", ErrGodotNotFound
	}

	// Define the target 'Executable'.
	ex := executable.New(v, p)

	if !store.Has(storePath, ex) {
		return "", ErrGodotNotFound
	}

	return store.ToolPath(storePath, ex)
}
