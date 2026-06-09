package scan

import (
	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func scanSectionHeaders(
	context runner.Context,
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
	context runner.Context,
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
	context runner.Context,
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
