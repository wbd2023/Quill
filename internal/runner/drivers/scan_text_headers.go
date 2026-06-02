package drivers

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
)

func runSectionHeadersScanner(
	context runner.Context,
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextConfig(context)
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

func runSectionHeaderDensityScanner(
	context runner.Context,
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextConfig(context)
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

func runSectionHeaderNamesScanner(
	context runner.Context,
) (result contract.ExecutionResult, err error) {
	config, err := decodeTextConfig(context)
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

func decodeTextConfig(context runner.Context) (config text.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(builtin.PackText)
	if !found {
		return text.Config{}, errMissingPackConfig(builtin.PackText)
	}

	return text.DecodeConfig(pack)
}
