package execute

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func Execute() *cli.Command {
	return &cli.Command{
		Name:  "execute",
		Usage: "execute a binary",
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
	input := c.String("input")

	if input == "" {
		return fmt.Errorf("input is required")
	}

	binary, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	// TODO: execute the binary
	fmt.Println("Executing binary...")
	fmt.Println("Input file:", input)
	fmt.Println("Binary:", binary)

	return nil
}
