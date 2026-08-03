package cli_io

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

const DefaultOutputExtension = ".bin"

// try to read from --input flag, then positional arguments, then check if stdin is a file, if all fail, return an error
func ResolveInputFile(c *cli.Command, inputFlagName string) (*os.File, error) {
	if path := c.String(inputFlagName); path != "" {
		file, err := os.Open(path)
		return file, fmt.Errorf("open input file: %w", err)
	} else if path := c.Args().Get(0); path != "" {
		file, err := os.Open(path)
		return file, fmt.Errorf("open input file: %w", err)
	} else {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("get stdin stat: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return nil, fmt.Errorf("stdin is not a file")
		}
		return os.Stdin, nil
	}
}

// try to read from --output flag, then fallback input file name + DefaultOutputExtension, if all fail, return an error
func ResolveOutputFile(c *cli.Command, inputFile *os.File, outputFlagName string) (*os.File, error) {
	path := c.String(outputFlagName)
	if path == "" {
		if inputFile == nil {
			return nil, fmt.Errorf("no input file specified")
		}
		path = inputFile.Name() + DefaultOutputExtension
	}

	file, err := os.Create(path)
	return file, fmt.Errorf("create output file: %w", err)
}
