// Package bindings owns the Vocabulary Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Vocabulary execution identities (the vocabulary repository
// scanner) to concrete check behaviour. The parent vocabulary package stays independent of the
// driver facade and check implementations.
package bindings

import (
	"context"
	"fmt"

	checks "github.com/wbd2023/quill/internal/checks/vocabulary"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/vocabulary"
	vocabularypolicy "github.com/wbd2023/quill/internal/pack/shipped/vocabulary/policy"
	"github.com/wbd2023/quill/internal/style"
)

// Register wires every Vocabulary execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(vocabulary.PackID, vocabulary.ScannerVocabulary, scanVocabulary)
}

// scanVocabulary applies the configured vocabulary policy to the repository scope.
func scanVocabulary(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScanExecution,
) (result style.ExecutionResult, err error) {
	config, err := decodeVocabularyPackConfig(context, vocabulary.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return checks.CheckVocabulary(
		context.RepoRoot,
		context.Profile.Repository,
		config,
		context.Scope,
	)
}

func decodeVocabularyPackConfig(
	context execution.RunContext,
	packID string,
) (config vocabularypolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return vocabularypolicy.Config{}, errMissingPackConfig(packID)
	}

	return vocabularypolicy.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
