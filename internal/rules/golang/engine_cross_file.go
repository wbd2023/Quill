package golang

import "ciphera/tools/internal/rules/golang/checks"

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	if state.enabled(GoCheckDomainIdentifiers) {
		typeAwareViolations, typeAwareRan := checks.CollectTypeAwareDomainIdentifierCastViolations(
			scanRoots,
			state.scannedGoFiles,
			state.pathClassifier,
			state.domainIdentifierConstructors,
		)
		if typeAwareRan {
			state.violations = append(state.violations, typeAwareViolations...)
		}
	}

	if state.collectOrder() {
		state.violations = append(state.violations, state.orderCollector.Violations()...)
	}
}
