package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/style"
)

/* ---------------------------------------- Check Command --------------------------------------- */

type checkCmd struct {
	repoFlags

	Scope style.Scope `kong:"name=scope,help=configured scope (profile default when omitted)"`

	Mode style.CheckMode `kong:"default=required,enum='required,all',help=mode: required|all"`

	Strict bool `kong:"name=strict-recommendations,help=fail on recommendation findings"`

	Verbose bool `kong:"name=verbose,help=print failing output"`
}

func (c *checkCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engineInstance, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("check", c.Format, err)
	}

	result, err := engineInstance.Check(ctx, engine.CheckOptions{
		Scope:                 c.Scope,
		Mode:                  c.Mode,
		StrictRecommendations: c.Strict,
	})
	if err != nil {
		return runner.reportCommandError("check", c.Format, err)
	}

	checkResult := report.NewCheckResult(result)

	summary, err := report.WriteCheck(
		runner.stdout,
		runner.envelopeMetadata("check"),
		c.Format,
		report.NewCheckView(checkResult),
		c.Verbose,
	)
	if err != nil {
		return runner.reportCommandError("check", c.Format, err)
	}

	if summary.Failed > 0 || summary.Blocked > 0 || summary.Errored > 0 {
		return 1
	}

	return 0
}
