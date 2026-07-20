package engine

import (
	"context"

	"github.com/wbd2023/Quill/internal/coverage"
	"github.com/wbd2023/Quill/internal/styleguide"
)

// Coverage loads STYLE.md and the compiled effective profile, then builds requirement coverage.
//
// Coverage intentionally does not construct a runner context or inspect tools. It uses a separate
// internal pipeline that loads only the profile and style guide.
func (engine *Engine) Coverage(
	operationContext context.Context,
) (coverageReport coverage.Report, operationError error) {
	compiled, err := engine.loadCompiledProfile(operationContext)
	if err != nil {
		return coverage.Report{}, err
	}

	document, err := styleguide.Load(engine.repositoryRoot, styleguide.Config{
		Filename: compiled.profile.Profile.StyleGuide.Path,
	})
	if err != nil {
		return coverage.Report{}, err
	}

	return coverage.Build(document, compiled.profile.Effective.Rules), nil
}
