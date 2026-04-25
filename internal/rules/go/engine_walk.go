package gostyle

import (
	"go/ast"
	"go/parser"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/rules/go/checks"
)

/* ------------------------------------------ File Walk ----------------------------------------- */

func (state *analysisState) walkDirectory(directory string) {
	walkError := filepath.WalkDir(
		directory,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if state.shouldSkipDirectory(entry) {
				return filepath.SkipDir
			}

			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			state.processFile(path)
			return nil
		},
	)
	if walkError != nil {
		state.writeWarning("error walking %s: %v\n", directory, walkError)
	}
}

func (state *analysisState) shouldSkipDirectory(entry os.DirEntry) (found bool) {
	if !entry.IsDir() {
		return false
	}

	for _, exclusion := range state.repository.GlobalExclusions {
		if entry.Name() == exclusion {
			return true
		}
	}

	return false
}

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
	state.violations = append(
		state.violations,
		checks.CheckProcessExecutionSafety(state.fileSet, file)...,
	)
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
	state.violations = append(state.violations, checks.CheckNamedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checks.CheckNakedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checks.CheckTypeElision(state.fileSet, file)...)
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
	state.violations = append(
		state.violations,
		checks.CheckInlineCommentStyle(
			state.fileSet,
			file,
			normalisedPath,
			state.pathClassifier,
		)...,
	)
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
	state.violations = append(
		state.violations,
		checks.CheckParameterOrder(state.fileSet, file, state.goParameters)...,
	)
	state.violations = append(
		state.violations,
		checks.CheckConstructorOrder(state.fileSet, file, state.goParameters)...,
	)
	if !isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckFileStructureOrder(state.fileSet, file)...,
		)
	}

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

	// Single-letter checks skip test files to reduce noise in table-driven structures and assertion
	// helpers.
	if !isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckSingleLetterVars(state.fileSet, file)...,
		)
	}

	if isTestFile {
		state.violations = append(
			state.violations,
			checks.CheckTestHygiene(state.fileSet, file, normalisedPath)...,
		)
	}
}

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	typeAwareViolations, typeAwareRan := checks.CollectTypeAwareDomainIdentifierCastViolations(
		scanRoots,
		state.scannedGoFiles,
		state.pathClassifier,
		state.goIdentifiers,
	)
	if typeAwareRan {
		state.violations = append(state.violations, typeAwareViolations...)
	}

	state.violations = append(state.violations, state.orderCollector.Violations()...)
}
