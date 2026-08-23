package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/report"
)

/* ---------------------------------------- List Command ---------------------------------------- */

type listCmd struct {
	repoFlags

	Selector string `kong:"arg,required,enum='packs,rules,tools,scopes',help=selector"`
}

func (c *listCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engineInstance, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("list", c.Format, err)
	}

	snapshot, err := engineInstance.Metadata(ctx)
	if err != nil {
		return runner.reportCommandError("list", c.Format, err)
	}

	result := report.NewListResult(snapshot, c.Selector)
	if err := report.WriteList(
		runner.stdout, runner.envelopeMetadata("list"), c.Format, result,
	); err != nil {
		return runner.reportCommandError("list", c.Format, err)
	}

	return 0
}
