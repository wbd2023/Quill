package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
)

func scanSectionHeaders(
	context runner.Context,
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
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
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
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
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextPackConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	return text.CheckSectionHeaderNames(
		context.RepoRoot,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}
