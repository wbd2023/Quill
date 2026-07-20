package main

import (
	"os"

	"github.com/wbd2023/Quill/internal/cli"
)

func main() {
	tool := cli.New(os.Stdout, os.Stderr, currentVersion())
	os.Exit(tool.Run(os.Args[1:]))
}
