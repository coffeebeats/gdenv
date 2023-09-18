package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

const (
	EnvGDEnvArch     = "GDENV_ARCH"
	EnvGDEnvOS       = "GDENV_OS"
	EnvGDEnvPlatform = "GDENV_PLATFORM"
)

// A 'urfave/cli' command to download and cache a specific version of Godot.
func NewInstall() *cli.Command {
	return &cli.Command{
		Name:     "install",
		Category: "Install",

		Usage:     "download and cache a specific version of Godot",
		UsageText: "gdenv install [OPTIONS] <VERSION>",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil {
				return failWithUsage(c, err)
			}

			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return fail(err)
			}

			// Define the host 'Platform'.
			p, err := detectPlatform()
			if err != nil {
				return fail(err)
			}

			// Define the target 'Executable'.
			ex := executable.New(v, p)

			if store.Has(storePath, ex) && !c.Bool("force") {
				return nil
			}

			if err := install(storePath, ex); err != nil {
				return fail(err)
			}

			return nil
		},
	}
}

/* ---------------------------- Function: install --------------------------- */

// Downloads and caches a platform-specific version of Godot.
func install(_ string, _ executable.Executable) error {
	return nil
}

/* ------------------------ Function: detectPlatform ------------------------ */

// Resolves the target platform by first checking environment variables and then
// falling back to the host platform.
func detectPlatform() (platform.Platform, error) {
	// First, check the full platform override.
	if platformRaw := os.Getenv(EnvGDEnvPlatform); platformRaw != "" {
		p, err := platform.Parse(platformRaw)
		if err != nil {
			return p, fmt.Errorf("%w: '%s'", err, platformRaw)
		}

		return p, nil
	}

	// Next, check the individual platform components for overrides and assemble
	// them into a 'Platform'.

	osRaw := os.Getenv(EnvGDEnvOS)
	if osRaw == "" {
		osRaw = runtime.GOOS
	}

	o, err := platform.ParseOS(osRaw)
	if err != nil {
		return platform.Platform{}, fmt.Errorf("%w: '%s'", err, osRaw)
	}

	archRaw := os.Getenv(EnvGDEnvArch)
	if archRaw == "" {
		archRaw = runtime.GOARCH
	}

	a, err := platform.ParseArch(archRaw)
	if err != nil {
		return platform.Platform{}, fmt.Errorf("%w: '%s'", err, archRaw)
	}

	return platform.Platform{Arch: a, OS: o}, nil
}
