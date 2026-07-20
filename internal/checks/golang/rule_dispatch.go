package golang

import (
	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
	"github.com/wbd2023/Quill/internal/checks/golang/check"
	"github.com/wbd2023/Quill/internal/checks/golang/structure"
	"github.com/wbd2023/Quill/internal/checks/golang/syntax"
	"github.com/wbd2023/Quill/internal/checks/golang/test"
)

/* ---------------------------------------- Rule Dispatch --------------------------------------- */

func (scan fileScan) addLoggingViolations() {
	if !scan.enabled(check.Logging) {
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
	if !scan.enabled(check.Security) {
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
	if !scan.enabled(check.Process) {
		return
	}

	scan.addViolations(syntax.CheckProcessExecutionSafety(
		scan.state.fileSet,
		scan.file,
	))
}

func (scan fileScan) addResourceViolations() {
	if !scan.enabled(check.Resources) {
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
	if !scan.enabled(check.Data) {
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
	if !scan.enabled(check.Returns) {
		return
	}

	scan.addViolations(syntax.CheckNamedReturns(scan.state.fileSet, scan.file))
	scan.addViolations(syntax.CheckNakedReturns(scan.state.fileSet, scan.file))
}

func (scan fileScan) addParameterViolations() {
	if !scan.enabled(check.Parameters) {
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
	if !scan.enabled(check.Errors) {
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
	if !scan.enabled(check.Comments) {
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
	if !scan.enabled(check.DomainValues) {
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
	if !scan.enabled(check.Order) {
		return
	}

	if !scan.isTestFile {
		scan.addViolations(structure.CheckStructureOrder(
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
	if !scan.enabled(check.Naming) || scan.isTestFile {
		return
	}

	scan.addViolations(syntax.CheckSingleLetterVars(scan.state.fileSet, scan.file))
}

func (scan fileScan) addTestViolations() {
	if !scan.enabled(check.Tests) || !scan.isTestFile {
		return
	}

	scan.addViolations(test.CheckHygiene(
		scan.state.fileSet,
		scan.file,
		scan.path,
	))
}

func (scan fileScan) addFileShapeViolations() {
	if !scan.enabled(check.FileShape) {
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
	if scan.enabled(check.GuardClauseSpacing) {
		scan.addViolations(structure.CheckGuardClauseSpacing(
			scan.state.fileSet,
			scan.file,
		))
	}

	if scan.enabled(check.SwitchCaseSpacing) {
		scan.addViolations(structure.CheckSwitchCaseSpacing(
			scan.state.fileSet,
			scan.file,
			scan.lines,
		))
	}
}

/* -------------------------------------- Dispatch Helpers -------------------------------------- */

func (scan fileScan) enabled(checkName string) (enabled bool) {
	return scan.state.enabled(checkName)
}

func (scan fileScan) addViolations(violations []analysis.Violation) {
	scan.state.violations = append(scan.state.violations, violations...)
}
