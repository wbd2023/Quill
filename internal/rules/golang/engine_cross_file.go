package golang

import (
	"ciphera/tools/internal/rules/golang/check"
	"ciphera/tools/internal/rules/golang/syntax"
)

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	if state.enabled(check.DomainValues) {
		typeAwareViolations, typeAwareRan := syntax.CollectTypeAwareDomainValueCastViolations(
			scanRoots,
			state.scannedGoFiles,
			state.pathClassifier,
			state.domainValueConstructors,
		)
		if typeAwareRan {
			state.violations = append(state.violations, typeAwareViolations...)
		}
	}

	if state.collectOrder() {
		state.violations = append(state.violations, state.orderCollector.Violations()...)
	}
}
