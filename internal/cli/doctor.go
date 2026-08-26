package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/report"
)

type doctorCmd struct {
	repoFlags
}

func (c *doctorCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engine, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("doctor", c.Format, err)
	}

	inspection, err := engine.Inspect(ctx)
	if err != nil {
		return runner.reportCommandError("doctor", c.Format, err)
	}

	allValid, err := report.WriteToolchain(
		runner.stdout,
		runner.envelopeMetadata("doctor"),
		c.Format,
		report.NewToolchainResult(inspection),
	)
	if err != nil {
		return runner.reportCommandError("doctor", c.Format, err)
	}

	if !allValid {
		return 1
	}

	return 0
}
