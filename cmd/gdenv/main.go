package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
)

func main() { //nolint:funlen
	log.SetFlags(0)

	cli.VersionPrinter = versionPrinter

	app := &cli.App{
		Name:    "gdenv",
		Version: "v0.4.0", // x-release-please-version

		Suggest:                true,
		UseShortOptionHandling: true,

		Commands: []*cli.Command{

			/* -------------------------------- Pin/Unpin ------------------------------- */

			NewPin(),
			NewUnpin(),

			/* ---------------------------- Install/Uninstall --------------------------- */

			NewInstall(),
			NewUninstall(),

			/* --------------------------------- Utility -------------------------------- */

			NewCompletions(),
			NewLs(),
			NewWhich(),
		},
	}

	// Call 'os.Exit' as the first-in/last-out defer; ensures an exit code is
	// returned to the caller.
	var exitCode int
	defer func() {
		if err := recover(); err != nil {
			exitCode = 1

			log.Println(err)
		}

		os.Exit(exitCode)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Ensure that the signal handler is removed after first interrupt.
	go func() {
		<-ctx.Done()
		stop()
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		var usageErr UsageError
		if errors.As(err, &usageErr) {
			usageErr.PrintUsage()
		}

		log.Println(err)

		exitCode = 1
	}
}

/* -------------------------------------------------------------------------- */
/*                              Type: UsageError                              */
/* -------------------------------------------------------------------------- */

// UsageError is any error returned from a subcommand implementation that should
// have subcommand usage instructions printed.
type UsageError struct {
	ctx *cli.Context
	err error
}

/* -------------------------- Function: PrintUsage -------------------------- */

// PrintUsage prints the usage associated with the subcommand that failed.
func (e UsageError) PrintUsage() {
	// NOTE: This never returns a meaningful error so ignore it.
	cli.ShowSubcommandHelp(e.ctx) //nolint:errcheck
}

/* ------------------------------- Impl: Error ------------------------------ */

func (e UsageError) Error() string {
	return e.err.Error()
}

/* -------------------------------------------------------------------------- */
/*                          Function: versionPrinter                          */
/* -------------------------------------------------------------------------- */

func versionPrinter(cCtx *cli.Context) {
	log.Printf("gdenv %s\n", cCtx.App.Version)
}
