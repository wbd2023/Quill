package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/checks/vocabulary"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
)

// CheckVocabulary check vocabulary.
func CheckVocabulary(vocabularyPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanVocabulary(ctx, context, execution, vocabularyPackID)
	}
}

func scanVocabulary(
	ctx context.Context,
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
