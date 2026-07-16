package scan

import (
	"ciphera/tools/internal/checks/vocabulary"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckVocabulary check vocabulary.
func CheckVocabulary(vocabularyPackID string) (scanner runtimebinding.RepositoryScanner) {
	return func(
		context execution.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanVocabulary(context, execution, vocabularyPackID)
	}
}

func scanVocabulary(
	context execution.Context,
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
