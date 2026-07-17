package scan

import (
	"ciphera/tools/internal/checks/vocabulary"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckVocabulary check vocabulary.
func CheckVocabulary(vocabularyPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanVocabulary(context, execution, vocabularyPackID)
	}
}

func scanVocabulary(
	context execution.RunContext,
	_ style.RepositoryScanExecution,
	vocabularyPackID string,
) (result style.ExecutionResult, err error) {
	config, err := decodeVocabularyPackConfig(context, vocabularyPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return vocabulary.CheckVocabulary(
		context.RepoRoot,
		context.Profile.Repository,
		config,
		context.Scope,
	)
}
