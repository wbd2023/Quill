package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/report"
)

type coverageCmd struct {
	repoFlags
	Verbose bool `kong:"name=verbose,help=print requirement-level detail"`
}

func (c *coverageCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	engine, err := c.newEngine()
	if err != nil {
		return runner.reportCommandError("coverage", c.Format, err)
	}

	coverage, err := engine.Coverage(ctx)
	if err != nil {
		return runner.reportCommandError("coverage", c.Format, err)
	}

	if err = report.WriteCoverage(
		runner.stdout,
		runner.envelopeMetadata("coverage"),
		c.Format,
		report.NewCoverageView(coverage),
		c.Verbose,
	); err != nil {
		return runner.reportCommandError("coverage", c.Format, err)
	}

	return 0
}
