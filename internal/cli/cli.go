package cli

import "io"

func New(stdout io.Writer, stderr io.Writer) (tool CLI) {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	return CLI{
		stdout:          stdout,
		stderr:          stderr,
		resolveRepoRoot: resolveRepoRoot,
	}
}
