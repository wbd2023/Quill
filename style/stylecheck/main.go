package main

import (
	"os"

	"stylecheck/internal/lint"
)

func main() {
	os.Exit(lint.Check(os.Args[1:]))
}
