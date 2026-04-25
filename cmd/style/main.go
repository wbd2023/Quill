package main

import (
	"os"

	"ciphera/tools/internal/cli"
)

func main() {
	tool := cli.New(os.Stdout, os.Stderr)
	os.Exit(tool.Run(os.Args[1:]))
}
