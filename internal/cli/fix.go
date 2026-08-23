package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/style"
)

/* ----------------------------------------- Fix Command ---------------------------------------- */

type fixCmd struct {
	repoFlags
	Scope style.Scope `kong:"name=scope,help=configured scope (profile default when omitted)"`
}

func (c *fixCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engineInstance, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("fix", c.Format, err)
	}

	result, err := engineInstance.Fix(ctx, engine.FixOptions{Scope: c.Scope})
	if err != nil {
		return runner.reportCommandError("fix", c.Format, err)
	}

	writer := runner.stdout
	if c.Format == report.FormatText {
		writer = runner.stderr
	}
	summary, err := report.WriteFix(
		writer,
		runner.envelopeMetadata("fix"),
		c.Format,
		report.NewFixResult(result),
	)
	if err != nil {
		return runner.reportCommandError("fix", c.Format, err)
	}
	if !summary.AllValid || summary.HasExecutionError {
		return 1
	}

	return 0
}
