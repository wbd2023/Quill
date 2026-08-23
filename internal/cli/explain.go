package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/report"
)

/* --------------------------------------- Explain Command -------------------------------------- */

type explainCmd struct {
	repoFlags
	Subject ruleRef `kong:"arg,required,help=subject: rule:<id>"`
}

func (c *explainCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engineInstance, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("explain", c.Format, err)
	}

	explanation, err := engineInstance.Explain(ctx, c.Subject.ruleID())
	if err != nil {
		return runner.reportCommandError("explain", c.Format, err)
	}

	if err := report.WriteExplain(
		runner.stdout,
		runner.envelopeMetadata("explain"),
		c.Format,
		report.NewExplainResult(explanation),
	); err != nil {
		return runner.reportCommandError("explain", c.Format, err)
	}

	return 0
}
