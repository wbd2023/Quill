package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/vocabulary"
	"ciphera/tools/internal/runner"
)

func vocabularyPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerVocabulary: scanVocabulary,
	}
}

func scanVocabulary(
	context runner.Context,
	_ contract.RepositoryScanExecution,
) (result contract.ExecutionResult, err error) {
	config, err := decodeVocabularyPackConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	return vocabulary.CheckVocabulary(
		context.RepoRoot,
		context.Profile.Repository,
		config,
		context.Scope,
	)
}
