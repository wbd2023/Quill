package golang

import "ciphera/tools/internal/rules/golang/checks"

/* ---------------------------------------- Rule Dispatch --------------------------------------- */

func (analysis fileAnalysis) addLoggingViolations() {
	if !analysis.enabled(GoCheckLogging) {
		return
	}

	analysis.addViolations(checks.CheckStructuredLogging(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.state.pathClassifier,
		analysis.state.goParameters,
	))
}

func (analysis fileAnalysis) addSecurityViolations() {
	if !analysis.enabled(GoCheckSecurity) {
		return
	}

	analysis.addViolations(checks.CheckSensitiveDataLiterals(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
		analysis.state.goParameters,
	))
	analysis.addViolations(checks.CheckCryptographySafety(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
	))
}

func (analysis fileAnalysis) addProcessViolations() {
	if !analysis.enabled(GoCheckProcess) {
		return
	}

	analysis.addViolations(checks.CheckProcessExecutionSafety(
		analysis.state.fileSet,
		analysis.file,
	))
}

func (analysis fileAnalysis) addResourceViolations() {
	if !analysis.enabled(GoCheckResources) {
		return
	}

	analysis.addViolations(checks.CheckContextAndResourceSafety(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
	))
}

func (analysis fileAnalysis) addDataViolations() {
	if !analysis.enabled(GoCheckData) {
		return
	}

	analysis.addViolations(checks.CheckDataUsage(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
	))
}

func (analysis fileAnalysis) addReturnViolations() {
	if !analysis.enabled(GoCheckReturns) {
		return
	}

	analysis.addViolations(checks.CheckNamedReturns(analysis.state.fileSet, analysis.file))
	analysis.addViolations(checks.CheckNakedReturns(analysis.state.fileSet, analysis.file))
}

func (analysis fileAnalysis) addParameterViolations() {
	if !analysis.enabled(GoCheckParameters) {
		return
	}

	analysis.addViolations(checks.CheckTypeElision(analysis.state.fileSet, analysis.file))
	analysis.addViolations(checks.CheckParameterOrder(
		analysis.state.fileSet,
		analysis.file,
		analysis.state.goParameters,
	))
	analysis.addViolations(checks.CheckConstructorOrder(
		analysis.state.fileSet,
		analysis.file,
		analysis.state.goParameters,
	))
}

func (analysis fileAnalysis) addErrorViolations() {
	if !analysis.enabled(GoCheckErrors) {
		return
	}

	analysis.addViolations(checks.CheckErrorHandlingStyle(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
		analysis.state.goParameters,
	))
	analysis.addViolations(checks.CheckAdapterErrorWrapping(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
		analysis.state.pathClassifier,
	))
}

func (analysis fileAnalysis) addCommentViolations() {
	if !analysis.enabled(GoCheckComments) {
		return
	}

	analysis.addViolations(checks.CheckInlineCommentStyle(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.state.pathClassifier,
	))
}

func (analysis fileAnalysis) addDomainIdentifierViolations() {
	if !analysis.enabled(GoCheckDomainIdentifiers) {
		return
	}

	analysis.addViolations(checks.CheckDirectDomainIdentifierCasts(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.state.pathClassifier,
		analysis.state.goIdentifiers,
	))
}

func (analysis fileAnalysis) addOrderViolations() {
	if !analysis.enabled(GoCheckOrder) {
		return
	}

	if !analysis.isTestFile {
		analysis.addViolations(checks.CheckFileStructureOrder(
			analysis.state.fileSet,
			analysis.file,
		))
		analysis.addViolations(checks.CheckScannerEntrypointOrder(
			analysis.state.fileSet,
			analysis.file,
			analysis.path,
		))
	}

	analysis.addViolations(checks.CheckCRUDLOrder(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.state.pathClassifier,
	))
	analysis.state.orderCollector.Collect(analysis.state.fileSet, analysis.file, analysis.path)
}

func (analysis fileAnalysis) addNamingViolations() {
	if !analysis.enabled(GoCheckNaming) || analysis.isTestFile {
		return
	}

	analysis.addViolations(checks.CheckSingleLetterVars(analysis.state.fileSet, analysis.file))
}

func (analysis fileAnalysis) addTestViolations() {
	if !analysis.enabled(GoCheckTests) || !analysis.isTestFile {
		return
	}

	analysis.addViolations(checks.CheckTestHygiene(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
	))
}

func (analysis fileAnalysis) addFileShapeViolations() {
	if !analysis.enabled(GoCheckFileShape) {
		return
	}

	analysis.addViolations(checks.CheckFileShape(
		analysis.state.fileSet,
		analysis.file,
		analysis.path,
		analysis.isTestFile,
	))
}

/* -------------------------------------- Dispatch Helpers -------------------------------------- */

func (analysis fileAnalysis) enabled(checkName string) (enabled bool) {
	return analysis.state.enabled(checkName)
}

func (analysis fileAnalysis) addViolations(violations []checks.Violation) {
	analysis.state.violations = append(analysis.state.violations, violations...)
}
