package golang

import (
	"fmt"
	"go/token"
	"io"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rules/golang/checks"
	"ciphera/tools/internal/rules/golang/order"
)

type analysisState struct {
	repository     policy.RepositoryConfig
	goParameters   policy.GoParameterConfig
	goIdentifiers  policy.GoDomainIdentifierConfig
	enabledChecks  map[string]bool
	pathClassifier checks.PathClassifier
	fileSet        *token.FileSet
	scannedGoFiles []string
	violations     []checks.Violation
	warningWriter  io.Writer
	orderCollector *order.Collector
}

func newAnalysisState(
	repoRoot string,
	config policy.Config,
	checkNames []string,
) (state *analysisState) {
	pathClassifier := checks.NewPathClassifier(repoRoot, config.Paths)

	return &analysisState{
		repository:     config.Repository,
		goParameters:   config.Naming.GoParameters,
		goIdentifiers:  config.Naming.GoDomainIdentifiers,
		enabledChecks:  enabledGoChecks(checkNames),
		pathClassifier: pathClassifier,
		fileSet:        token.NewFileSet(),
		scannedGoFiles: make([]string, 0),
		warningWriter:  io.Discard,
		orderCollector: order.NewCollector(pathClassifier),
	}
}

func enabledGoChecks(checkNames []string) (enabled map[string]bool) {
	enabled = make(map[string]bool, len(checkNames))
	for _, checkName := range checkNames {
		enabled[checkName] = true
	}

	return enabled
}

func (state *analysisState) enabled(checkName string) (enabled bool) {
	if len(state.enabledChecks) == 0 {
		return true
	}

	return state.enabledChecks[checkName]
}

func (state *analysisState) collectOrder() (collect bool) {
	return state.enabled(GoCheckOrder)
}

func (state *analysisState) writeWarning(format string, arguments ...any) {
	if state.warningWriter == nil {
		return
	}

	_, _ = fmt.Fprintf(state.warningWriter, format, arguments...)
}
