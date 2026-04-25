package gostyle

import (
	"fmt"
	"go/token"
	"io"

	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rules/go/checks"
	ruleorder "ciphera/tools/internal/rules/go/order"
)

/* --------------------------------------- Analysis State --------------------------------------- */

type analysisState struct {
	repository     profile.RepositoryConfig
	goParameters   profile.GoParameterConfig
	goIdentifiers  profile.GoDomainIdentifierConfig
	pathClassifier checks.PathClassifier
	fileSet        *token.FileSet
	scannedGoFiles []string
	violations     []checks.Violation
	warningWriter  io.Writer
	orderCollector *ruleorder.Collector
}

func newAnalysisState(
	repoRoot string,
	policy profile.Profile,
) (state *analysisState) {
	pathClassifier := checks.NewPathClassifier(repoRoot, policy.Paths)

	return &analysisState{
		repository:     policy.Repository,
		goParameters:   policy.Naming.GoParameters,
		goIdentifiers:  policy.Naming.GoDomainIdentifiers,
		pathClassifier: pathClassifier,
		fileSet:        token.NewFileSet(),
		scannedGoFiles: make([]string, 0),
		warningWriter:  io.Discard,
		orderCollector: ruleorder.NewCollector(pathClassifier),
	}
}

func (state *analysisState) writeWarning(format string, arguments ...any) {
	if state.warningWriter == nil {
		return
	}

	_, _ = fmt.Fprintf(state.warningWriter, format, arguments...)
}
