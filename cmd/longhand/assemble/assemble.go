package assemble

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func Assemble() *cli.Command {
	return &cli.Command{
		Name:  "assemble",
		Usage: "assemble a binary",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "output file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "disassemble",
				Usage: "disassemble mode",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "debug mode",
				Value: false,
			},
		},
		Action: assembleAction,
	}
}

func assembleAction(ctx context.Context, c *cli.Command) error {
	input := c.String("input")
	output := c.String("output")

	if input == "" {
		return fmt.Errorf("input is required")
	}
	if output == "" {
		return fmt.Errorf("output is required")
	}

	binary, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	// TODO: assemble the binary
	fmt.Println("Assembling binary...")
	fmt.Println("Input file:", input)
	fmt.Println("Output file:", output)
	fmt.Println("Binary:", binary)

	err = os.WriteFile(output, binary, 0644)
	if err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}
