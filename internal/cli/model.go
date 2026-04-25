package cli

import "io"

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	helpCommand   = "help"
	usageExitCode = 2
)

/* -------------------------------------------- Types ------------------------------------------- */

type repositoryRootResolver func(string) (string, error)

type CLI struct {
	stdout          io.Writer
	stderr          io.Writer
	resolveRepoRoot repositoryRootResolver
}

type Action func(CLI) int
