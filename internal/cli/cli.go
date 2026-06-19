package cli

import "io"

func New(stdout io.Writer, stderr io.Writer) (tool Tool) {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	return Tool{
		stdout:          stdout,
		stderr:          stderr,
		resolveRepoRoot: resolveRepoRoot,
	}
}
