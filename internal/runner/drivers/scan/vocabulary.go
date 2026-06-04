package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/vocabulary"
	"ciphera/tools/internal/runner"
)

func vocabularyPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerVocabulary: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			pack, found := context.Profile.PackConfigs.Lookup(builtin.PackVocabulary)
			if !found {
				return contract.ExecutionResult{}, errMissingPackConfig(builtin.PackVocabulary)
			}

			config, err := vocabulary.DecodeConfig(pack)
			if err != nil {
				return contract.ExecutionResult{}, err
			}

			return vocabulary.CheckVocabulary(
				context.RepoRoot,
				context.Profile.Repository,
				config,
				context.Scope,
			)
		},
	}
}
