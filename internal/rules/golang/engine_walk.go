package golang

import (
	"go/ast"
	"go/parser"
	"strings"

	"ciphera/tools/internal/rules/golang/checks"
)

/* ----------------------------------------- File Checks ---------------------------------------- */

func (state *analysisState) processFile(path string) {
	file, parseError := parser.ParseFile(state.fileSet, path, nil, parser.ParseComments)
	if parseError != nil {
		state.writeWarning("warning: skipping %s: %v\n", path, parseError)
		return
	}

	normalisedPath := normalisePath(path)
	state.scannedGoFiles = append(state.scannedGoFiles, normalisedPath)
	isTestFile := strings.HasSuffix(path, "_test.go")
	state.addPerFileViolations(file, normalisedPath, isTestFile)
}

func (state *analysisState) addPerFileViolations(
	file *ast.File,
	normalisedPath string,
	isTestFile bool,
) {
	if state.enabled(GoCheckLogging) {
		state.violations = append(
			state.violations,
			checks.CheckStructuredLogging(
				state.fileSet,
				file,
				normalisedPath,
				state.pathClassifier,
				state.goParameters,
			)...,
		)
	}

	if state.enabled(GoCheckSecurity) {
		state.violations = append(
			state.violations,
			checks.CheckSensitiveDataLiterals(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
				state.goParameters,
			)...,
		)
		state.violations = append(
			state.violations,
			checks.CheckCryptographySafety(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
			)...,
		)
	}

	if state.enabled(GoCheckProcess) {
		state.violations = append(
			state.violations,
			checks.CheckProcessExecutionSafety(state.fileSet, file)...,
		)
	}

	if state.enabled(GoCheckResources) {
		state.violations = append(
			state.violations,
			checks.CheckContextAndResourceSafety(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
			)...,
		)
	}

	if state.enabled(GoCheckData) {
		state.violations = append(
			state.violations,
			checks.CheckDataUsage(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
			)...,
		)
	}

	if state.enabled(GoCheckReturns) {
		state.violations = append(
			state.violations,
			checks.CheckNamedReturns(state.fileSet, file)...,
		)
		state.violations = append(
			state.violations,
			checks.CheckNakedReturns(state.fileSet, file)...,
		)
	}

	if state.enabled(GoCheckParameters) {
		state.violations = append(state.violations, checks.CheckTypeElision(state.fileSet, file)...)
		state.violations = append(
			state.violations,
			checks.CheckParameterOrder(state.fileSet, file, state.goParameters)...,
		)
		state.violations = append(
			state.violations,
			checks.CheckConstructorOrder(state.fileSet, file, state.goParameters)...,
		)
	}

	if state.enabled(GoCheckErrors) {
		state.violations = append(
			state.violations,
			checks.CheckErrorHandlingStyle(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
				state.goParameters,
			)...,
		)
		state.violations = append(
			state.violations,
			checks.CheckAdapterErrorWrapping(
				state.fileSet,
				file,
				normalisedPath,
				isTestFile,
				state.pathClassifier,
			)...,
		)
	}

	if state.enabled(GoCheckComments) {
		state.violations = append(
			state.violations,
			checks.CheckInlineCommentStyle(
				state.fileSet,
				file,
				normalisedPath,
				state.pathClassifier,
			)...,
		)
	}

	if state.enabled(GoCheckDomainIdentifiers) {
		state.violations = append(
			state.violations,
			checks.CheckDirectDomainIdentifierCasts(
				state.fileSet,
				file,
				normalisedPath,
				state.pathClassifier,
				state.goIdentifiers,
			)...,
		)
	}

	if state.enabled(GoCheckOrder) && !isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckFileStructureOrder(state.fileSet, file)...,
		)
	}

	if state.enabled(GoCheckOrder) {
		state.violations = append(
			state.violations,
			checks.CheckCRUDLOrder(
				state.fileSet,
				file,
				normalisedPath,
				state.pathClassifier,
			)...,
		)

		state.orderCollector.Collect(state.fileSet, file, normalisedPath)
	}

	// Single-letter checks skip test files to reduce noise in table-driven structures and assertion
	// helpers.
	if state.enabled(GoCheckNaming) && !isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckSingleLetterVars(state.fileSet, file)...,
		)
	}

	if state.enabled(GoCheckTests) && isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckTestHygiene(state.fileSet, file, normalisedPath)...,
		)
	}
}

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	if state.enabled(GoCheckDomainIdentifiers) {
		typeAwareViolations, typeAwareRan := checks.CollectTypeAwareDomainIdentifierCastViolations(
			scanRoots,
			state.scannedGoFiles,
			state.pathClassifier,
			state.goIdentifiers,
		)
		if typeAwareRan {
			state.violations = append(state.violations, typeAwareViolations...)
		}
	}

	if state.collectOrder() {
		state.violations = append(state.violations, state.orderCollector.Violations()...)
	}
}
