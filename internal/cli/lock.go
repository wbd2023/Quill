package cli

import (
	"context"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

type lockCmd struct {
	repoFlags
}

func (c *lockCmd) run(ctx context.Context, runner Runner) (exitCode int) {
	progressWriter := runner.stdout
	if c.Format == report.FormatJSON {
		// Machine mode reserves stdout for the single envelope; route lock progress to stderr.
		progressWriter = runner.stderr
	}

	option := engine.WithProgressWriter(progressWriter)

	engine, err := c.newEngine(option)
	if err != nil {
		return runner.reportCommandError("lock", c.Format, err)
	}

	result, err := engine.Lock(ctx)
	if err != nil {
		return runner.reportCommandError("lock", c.Format, err)
	}

	if err := report.WriteLock(
		runner.stdout,
		runner.envelopeMetadata("lock"),
		c.Format,
		report.NewLockResult(result),
	); err != nil {
		return runner.reportCommandError("lock", c.Format, err)
	}

	return 0
}
