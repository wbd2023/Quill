package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/checks/text"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/style"
)

func scanSectionHeaders(
	ctx context.Context,
	context execution.RunContext,
	_ style.RepositoryScanExecution,
	textPackID string,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context, textPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckSectionHeaders(
		context.RepoRoot,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}

func scanSectionHeaderDensity(
	ctx context.Context,
	context execution.RunContext,
	_ style.RepositoryScanExecution,
	textPackID string,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context, textPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckSectionHeaderDensity(
		context.RepoRoot,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}

func scanSectionHeaderNames(
	ctx context.Context,
	context execution.RunContext,
	_ style.RepositoryScanExecution,
	textPackID string,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context, textPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckSectionHeaderNames(
		context.RepoRoot,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}
