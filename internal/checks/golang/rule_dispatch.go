package golang

import (
	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	"github.com/wbd2023/quill/internal/checks/golang/structure"
	"github.com/wbd2023/quill/internal/checks/golang/syntax"
	"github.com/wbd2023/quill/internal/checks/golang/test"
)

/* ---------------------------------------- Rule Dispatch --------------------------------------- */

func (scan fileScan) addLoggingViolations() {
	if !scan.enabled(CheckLogging) {
		return
	}

	scan.addViolations(syntax.CheckStructuredLogging(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.state.pathClassifier,
		scan.state.goParameters,
	))
}

func (scan fileScan) addSecurityViolations() {
	if !scan.enabled(CheckSecurity) {
		return
	}

	scan.addViolations(syntax.CheckSensitiveDataLiterals(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
		scan.state.goParameters,
	))
	scan.addViolations(syntax.CheckCryptographySafety(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
	))
}

func (scan fileScan) addProcessViolations() {
	if !scan.enabled(CheckProcess) {
		return
	}

	scan.addViolations(syntax.CheckProcessExecutionSafety(
		scan.state.fileSet,
		scan.file,
	))
}

func (scan fileScan) addResourceViolations() {
	if !scan.enabled(CheckResources) {
		return
	}

	scan.addViolations(syntax.CheckContextAndResourceSafety(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
	))
}

func (scan fileScan) addDataViolations() {
	if !scan.enabled(CheckData) {
		return
	}

	scan.addViolations(syntax.CheckDataUsage(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
	))
}

func (scan fileScan) addReturnViolations() {
	if !scan.enabled(CheckReturns) {
		return
	}

	scan.addViolations(syntax.CheckNamedReturns(scan.state.fileSet, scan.file))
	scan.addViolations(syntax.CheckNakedReturns(scan.state.fileSet, scan.file))
}

func (scan fileScan) addParameterViolations() {
	if !scan.enabled(CheckParameters) {
		return
	}

	scan.addViolations(syntax.CheckTypeElision(scan.state.fileSet, scan.file))
	scan.addViolations(syntax.CheckParameterOrder(
		scan.state.fileSet,
		scan.file,
		scan.state.goParameters,
	))
	scan.addViolations(syntax.CheckConstructorOrder(
		scan.state.fileSet,
		scan.file,
		scan.state.goConstructors,
		scan.state.goParameters,
	))
}

func (scan fileScan) addErrorViolations() {
	if !scan.enabled(CheckErrors) {
		return
	}

	scan.addViolations(syntax.CheckErrorHandlingStyle(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
		scan.state.goParameters,
	))
	scan.addViolations(syntax.CheckAdapterErrorWrapping(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
		scan.state.pathClassifier,
	))
}

func (scan fileScan) addCommentViolations() {
	if !scan.enabled(CheckComments) {
		return
	}

	scan.addViolations(syntax.CheckInlineCommentStyle(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.state.pathClassifier,
	))
}

func (scan fileScan) addDomainValueViolations() {
	if !scan.enabled(CheckDomainValues) {
		return
	}

	scan.addViolations(syntax.CheckDirectDomainValueCasts(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.state.pathClassifier,
		scan.state.domainValueConstructors,
	))
}

func (scan fileScan) addOrderViolations() {
	if !scan.enabled(CheckOrder) {
		return
	}

	if !scan.isTestFile {
		scan.addViolations(structure.CheckOrder(
			scan.state.fileSet,
			scan.file,
		))
		scan.addViolations(structure.CheckScannerEntrypointOrder(
			scan.state.fileSet,
			scan.file,
			scan.path,
		))
	}

	scan.addViolations(structure.CheckCRUDLOrder(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.state.pathClassifier,
	))
	scan.state.orderCollector.Collect(scan.state.fileSet, scan.file, scan.path)
}

func (scan fileScan) addNamingViolations() {
	if !scan.enabled(CheckNaming) || scan.isTestFile {
		return
	}

	scan.addViolations(syntax.CheckSingleLetterVars(scan.state.fileSet, scan.file))
	scan.addViolations(syntax.CheckPackageStutter(scan.state.fileSet, scan.file, scan.lines))
	scan.addViolations(syntax.CheckWeightlessSuffixes(scan.state.fileSet, scan.file, scan.lines))
}

func (scan fileScan) addTestViolations() {
	if !scan.enabled(CheckTests) || !scan.isTestFile {
		return
	}

	scan.addViolations(test.CheckHygiene(
		scan.state.fileSet,
		scan.file,
		scan.path,
	))
}

func (scan fileScan) addFileShapeViolations() {
	if !scan.enabled(CheckFileShape) {
		return
	}

	scan.addViolations(structure.CheckShape(
		scan.state.fileSet,
		scan.file,
		scan.path,
		scan.isTestFile,
	))
}

func (scan fileScan) addSpacingViolations() {
	if scan.enabled(CheckGuardClauseSpacing) {
		scan.addViolations(structure.CheckGuardClauseSpacing(
			scan.state.fileSet,
			scan.file,
		))
	}

	if scan.enabled(CheckSwitchCaseSpacing) {
		scan.addViolations(structure.CheckSwitchCaseSpacing(
			scan.state.fileSet,
			scan.file,
			scan.lines,
		))
	}
}

/* -------------------------------------- Dispatch Helpers -------------------------------------- */

func (scan fileScan) enabled(selector Check) (enabled bool) {
	return scan.state.enabled(selector)
}

func (scan fileScan) addViolations(violations []analysis.Violation) {
	scan.state.violations = append(scan.state.violations, violations...)
}
