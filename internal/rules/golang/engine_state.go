package golang

import (
	"fmt"
	"go/token"
	"io"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rules/golang/analysis"
	"ciphera/tools/internal/rules/golang/check"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
	"ciphera/tools/internal/rules/golang/relationships"
)

type analysisState struct {
	repository              policy.RepositoryConfig
	goParameters            gopolicy.ParameterConfig
	goConstructors          gopolicy.ConstructorConfig
	domainValueConstructors gopolicy.DomainValueConstructors
	enabledChecks           map[string]bool
	pathClassifier          analysis.PathClassifier
	fileSet                 *token.FileSet
	scannedGoFiles          []string
	violations              []analysis.Violation
	warningWriter           io.Writer
	orderCollector          *relationships.Collector
}

func newAnalysisState(
	repoRoot string,
	repository policy.RepositoryConfig,
	paths policy.PathRoles,
	goConfig gopolicy.Config,
	checkNames []string,
) (state *analysisState) {
	pathClassifier := analysis.NewPathClassifier(repoRoot, paths)

	return &analysisState{
		repository:              repository,
		goParameters:            goConfig.Parameters,
		goConstructors:          goConfig.Constructors,
		domainValueConstructors: goConfig.DomainValues.RequiredConstructors,
		enabledChecks:           enabledGoChecks(checkNames),
		pathClassifier:          pathClassifier,
		fileSet:                 token.NewFileSet(),
		scannedGoFiles:          make([]string, 0),
		warningWriter:           io.Discard,
		orderCollector:          relationships.NewCollector(pathClassifier),
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
	return state.enabled(check.Order)
}

func (state *analysisState) writeWarning(format string, arguments ...any) {
	if state.warningWriter == nil {
		return
	}

	_, _ = fmt.Fprintf(state.warningWriter, format, arguments...)
}
