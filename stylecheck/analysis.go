package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/* --------------------------------------- Analysis State --------------------------------------- */

func newAnalysisState() (state *analysisState) {
	return &analysisState{
		fileSet:                token.NewFileSet(),
		scannedGoFiles:         make([]string, 0),
		interfaces:             make(map[string]interfaceDecl),
		mocks:                  make(map[string][]methodDecl),
		implementations:        make(map[string][]methodDecl),
		implementationBindings: make([]implementationBinding, 0),
	}
}

/* ------------------------------------------ File Walk ----------------------------------------- */

func (state *analysisState) walkDirectory(directory string) {
	walkError := filepath.WalkDir(
		directory,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if shouldSkipDirectory(entry) {
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
		fmt.Fprintf(os.Stderr, "error walking %s: %v\n", directory, walkError)
	}
}

func shouldSkipDirectory(entry os.DirEntry) (found bool) {
	if !entry.IsDir() {
		return false
	}

	switch entry.Name() {
	case "vendor", ".git", "testdata":
		return true
	default:
		return false
	}
}

func (state *analysisState) processFile(path string) {
	file, parseError := parser.ParseFile(state.fileSet, path, nil, parser.ParseComments)
	if parseError != nil {
		fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", path, parseError)
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
	state.violations = append(state.violations, checkNamedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checkNakedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checkTypeElision(state.fileSet, file)...)
	state.violations = append(
		state.violations,
		checkGoErrorHandlingStyle(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		checkAdapterErrorWrapping(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		checkInlineCommentStyle(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		checkDirectDomainIdentifierCasts(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(state.violations, checkParamOrder(state.fileSet, file)...)
	state.violations = append(state.violations, checkConstructorOrder(state.fileSet, file)...)
	if !isTestFile {
		state.violations = append(state.violations, checkFileStructureOrder(state.fileSet, file)...)
	}

	state.violations = append(
		state.violations,
		checkServiceTypeNaming(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		checkCRUDLOrder(state.fileSet, file, normalisedPath)...,
	)

	collectInterfaces(state.fileSet, file, normalisedPath, state.interfaces)
	collectMockMethods(state.fileSet, file, normalisedPath, state.mocks)
	collectImplementationMethods(state.fileSet, file, normalisedPath, state.implementations)
	collectImplementationBindings(
		state.fileSet,
		file,
		normalisedPath,
		&state.implementationBindings,
	)

	// Single-letter checks skip test files to reduce noise in table-driven structures and assertion
	// helpers.
	if !isTestFile {
		state.violations = append(state.violations, checkSingleLetterVars(state.fileSet, file)...)
	}
}

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	typeAwareViolations, typeAwareRan := collectTypeAwareDomainIdentifierCastViolations(
		scanRoots,
		state.scannedGoFiles,
	)
	if typeAwareRan {
		state.violations = append(state.violations, typeAwareViolations...)
	}

	state.violations = append(
		state.violations,
		checkMockOrderAgainstInterfaces(state.interfaces, state.mocks)...,
	)
	state.violations = append(
		state.violations,
		checkImplementationOrderAgainstInterfaces(
			state.interfaces,
			state.implementations,
			state.implementationBindings,
		)...,
	)
}

/* ------------------------------------------ Reporting ----------------------------------------- */

func sortViolations(violations []violation) {
	sort.Slice(violations, func(i int, j int) bool {
		if violations[i].position.Filename == violations[j].position.Filename {
			return violations[i].position.Line < violations[j].position.Line
		}
		return violations[i].position.Filename < violations[j].position.Filename
	})
}

func dedupeViolations(violations []violation) (deduped []violation) {
	seen := make(map[string]bool)
	deduped = make([]violation, 0, len(violations))

	for _, current := range violations {
		key := fmt.Sprintf(
			"%s:%d:%d|%s|%s",
			current.position.Filename,
			current.position.Line,
			current.position.Column,
			current.rule,
			current.message,
		)

		if seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, current)
	}

	return deduped
}

func printViolationsAndExit(violations []violation) {
	if len(violations) == 0 {
		return
	}

	for _, current := range violations {
		fmt.Fprintf(os.Stderr, "%s: [%s] %s\n", current.position, current.rule, current.message)
	}
	os.Exit(1)
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}
