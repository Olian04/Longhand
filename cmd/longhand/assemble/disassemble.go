package assemble

import (
	"context"
	"fmt"

	"github.com/Olian04/Longhand/cmd/longhand/cli_io"
	"github.com/urfave/cli/v3"
)

func Disassemble() *cli.Command {
	return &cli.Command{
		Name:      "disassemble",
		Usage:     "disassemble a binary",
		ArgsUsage: "[arguments...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "input file",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "output file",
				Required: false,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "debug mode",
				Value: false,
			},
		},
		Action: disassembleAction,
	}
}

func disassembleAction(ctx context.Context, c *cli.Command) error {
	input, err := cli_io.ResolveInputFile(c, "input")
	if err != nil {
		return fmt.Errorf("resolve input file: %w", err)
	}
	defer input.Close()

	output, err := cli_io.ResolveOutputFile(c, input, "output")
	if err != nil {
		return fmt.Errorf("resolve output file: %w", err)
	}
	defer output.Close()

	_ = c.Bool("debug")

	// TODO: assemble the binary
	fmt.Println("Assembling binary...")
	fmt.Println("Input file:", input.Name())
	fmt.Println("Output file:", output.Name())
	return nil
}
