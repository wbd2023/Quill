package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

/* --------------------------------------- Install Command -------------------------------------- */

type installCmd struct {
	repoFlags
}

func (c *installCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	progressWriter := runner.stdout
	if c.Format == report.FormatJSON {
		// Machine mode reserves stdout for the single envelope; route install progress to stderr.
		progressWriter = runner.stderr
	}

	engineInstance, err := c.newEngine(engine.WithProgressWriter(progressWriter))
	if err != nil {
		return runner.reportCommandError("install", c.Format, err)
	}

	result, err := engineInstance.Install(ctx)
	if err != nil {
		return runner.reportCommandError("install", c.Format, err)
	}

	allValid, err := report.WriteInstall(
		runner.stdout,
		runner.envelopeMetadata("install"),
		c.Format,
		report.NewToolchainResult(result.Toolchain),
	)
	if err != nil {
		return runner.reportCommandError("install", c.Format, err)
	}
	if !allValid {
		return 1
	}

	return 0
}
