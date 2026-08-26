package golang

import "github.com/wbd2023/quill/internal/checks/golang/syntax"

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	if state.enabled(CheckDomainValues) {
		violations, ran := syntax.CollectTypeAwareDomainValueCastViolations(
			scanRoots,
			state.scannedGoFiles,
			state.pathClassifier,
			state.domainValueConstructors,
		)
		if ran {
			state.violations = append(state.violations, violations...)
		}
	}

	if state.collectOrder() {
		state.violations = append(state.violations, state.orderCollector.Violations()...)
	}
}
