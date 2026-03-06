package main

import (
	"os"

	"stylecheck/internal/checker"
)

func main() {
	os.Exit(checker.Run(os.Args[1:]))
}
