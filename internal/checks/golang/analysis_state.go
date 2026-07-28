package golang

import (
	"fmt"
	"go/token"
	"io"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	"github.com/wbd2023/quill/internal/checks/golang/relationships"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
	"github.com/wbd2023/quill/internal/policy"
)

type analysisState struct {
	repository              policy.RepositoryConfig
	goParameters            gopolicy.ParameterConfig
	goConstructors          gopolicy.ConstructorConfig
	domainValueConstructors gopolicy.DomainValueConstructors
	enabledChecks           map[Check]bool
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
	checks []Check,
) (state *analysisState) {
	pathClassifier := analysis.NewPathClassifier(repoRoot, paths)

	return &analysisState{
		repository:              repository,
		goParameters:            goConfig.Parameters,
		goConstructors:          goConfig.Constructors,
		domainValueConstructors: goConfig.DomainValues.RequiredConstructors,
		enabledChecks:           enabledGoChecks(checks),
		pathClassifier:          pathClassifier,
		fileSet:                 token.NewFileSet(),
		scannedGoFiles:          make([]string, 0),
		warningWriter:           io.Discard,
		orderCollector:          relationships.NewCollector(pathClassifier),
	}
}

func enabledGoChecks(checks []Check) (enabled map[Check]bool) {
	enabled = make(map[Check]bool, len(checks))
	for _, selector := range checks {
		enabled[selector] = true
	}

	return enabled
}

func (state *analysisState) enabled(selector Check) (enabled bool) {
	if len(state.enabledChecks) == 0 {
		return true
	}

	return state.enabledChecks[selector]
}

func (state *analysisState) collectOrder() (collect bool) {
	return state.enabled(CheckOrder)
}

func (state *analysisState) writeWarning(format string, arguments ...any) {
	if state.warningWriter == nil {
		return
	}

	_, _ = fmt.Fprintf(state.warningWriter, format, arguments...)
}
