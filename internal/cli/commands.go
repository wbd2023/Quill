package cli

import "io"

const (
	helpCommand   = "help"
	usageExitCode = 2
)

type repositoryRootResolver func(string) (string, error)

// Tool is tool.
type Tool struct {
	stdout          io.Writer
	stderr          io.Writer
	resolveRepoRoot repositoryRootResolver
}

// Action is action.
type Action func(Tool) int
