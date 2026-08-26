package engine

import (
	"context"

	"github.com/wbd2023/quill/internal/coverage"
)

// Coverage builds requirement coverage from the prepared STYLE.md document and the compiled Plan.
//
// Coverage is metadata-only: it shares the document loaded during operation preparation and never
// constructs a runner context, resolves drivers, or inspects tools.
func (engine *Engine) Coverage(
	ctx context.Context,
) (report coverage.Report, err error) {
	if err := ctx.Err(); err != nil {
		return coverage.Report{}, err
	}

	prepared, err := engine.prepare(ctx)
	if err != nil {
		return coverage.Report{}, err
	}

	if err := ctx.Err(); err != nil {
		return coverage.Report{}, err
	}

	return coverage.Build(prepared.document, prepared.plan.Rules), nil
}
