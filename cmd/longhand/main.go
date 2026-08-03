package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"

	"github.com/Olian04/Longhand/cmd/longhand/assemble"
	"github.com/Olian04/Longhand/cmd/longhand/execute"
	"github.com/Olian04/Longhand/cmd/longhand/version"
)

func main() {
	vi := version.Info()

	cli.VersionPrinter = printVersion(vi)

	root := &cli.Command{
		Name:    "longhand",
		Version: vi.Version,
		Commands: []*cli.Command{
			assemble.Assemble(),
			execute.Execute(),
		},
	}

	// SIGINT/SIGTERM cancel ctx so long-running commands can stop cleanly and
	// runCLI's deferred logging cleanup still executes.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := root.Run(ctx, os.Args); err != nil {
		// Written straight to stderr, not slog: runCLI's deferred cleanup has
		// already torn the configured logger down by the time Run returns.
		fmt.Fprintf(os.Stderr, "longhand: %v\n", err)
		os.Exit(1)
	}
}

func printVersion(vi version.VersionInfo) func(cmd *cli.Command) {
	return func(cmd *cli.Command) {
		_, err := fmt.Fprintf(cmd.Root().Writer, "%s version %s\nrevision %s\nbuild_time %s\n",
			cmd.Name, vi.Version, vi.Revision, vi.BuildTime)
		if err != nil {
			slog.Error("write version", "error_message", err.Error())
		}
	}
}
