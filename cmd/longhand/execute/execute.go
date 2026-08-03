package execute

import (
	"context"
	"fmt"

	"github.com/Olian04/Longhand/cmd/longhand/cli_io"
	"github.com/urfave/cli/v3"
)

func Execute() *cli.Command {
	return &cli.Command{
		Name:      "execute",
		Usage:     "execute a binary",
		ArgsUsage: "[arguments...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "input file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "debug mode",
				Value: false,
			},
		},
		Action: executeAction,
	}
}

func executeAction(ctx context.Context, c *cli.Command) error {
	input, err := cli_io.ResolveInputFile(c, "input")
	if err != nil {
		return fmt.Errorf("resolve input file: %w", err)
	}
	defer input.Close()

	_ = c.Bool("debug")

	// TODO: execute the binary
	fmt.Println("Executing binary...")
	fmt.Println("Input file:", input.Name())

	return nil
}
