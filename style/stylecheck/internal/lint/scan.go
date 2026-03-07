package lint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"stylecheck/internal/lint/index"
	"stylecheck/internal/lint/report"
	"stylecheck/internal/lint/rules"
)

/* --------------------------------------- Analysis State --------------------------------------- */

type analysisState struct {
	fileSet                *token.FileSet
	scannedGoFiles         []string
	violations             []report.Violation
	interfaces             map[string]index.InterfaceDecl
	mocks                  map[string][]index.MethodDecl
	implementations        map[string][]index.MethodDecl
	implementationBindings []index.ImplementationBinding
}

func newAnalysisState() (state *analysisState) {
	return &analysisState{
		fileSet:                token.NewFileSet(),
		scannedGoFiles:         make([]string, 0),
		interfaces:             make(map[string]index.InterfaceDecl),
		mocks:                  make(map[string][]index.MethodDecl),
		implementations:        make(map[string][]index.MethodDecl),
		implementationBindings: make([]index.ImplementationBinding, 0),
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
	state.violations = append(state.violations, rules.CheckNamedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, rules.CheckNakedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, rules.CheckTypeElision(state.fileSet, file)...)
	state.violations = append(
		state.violations,
		rules.CheckGoErrorHandlingStyle(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		rules.CheckAdapterErrorWrapping(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		rules.CheckInlineCommentStyle(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		rules.CheckDirectDomainIdentifierCasts(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(state.violations, rules.CheckParamOrder(state.fileSet, file)...)
	state.violations = append(state.violations, rules.CheckConstructorOrder(state.fileSet, file)...)
	if !isTestFile {
		state.violations = append(state.violations, rules.CheckFileStructureOrder(state.fileSet, file)...)
	}

	state.violations = append(
		state.violations,
		rules.CheckServiceTypeNaming(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		rules.CheckCRUDLOrder(state.fileSet, file, normalisedPath)...,
	)

	index.CollectInterfaces(state.fileSet, file, normalisedPath, state.interfaces)
	index.CollectMockMethods(state.fileSet, file, normalisedPath, state.mocks)
	index.CollectImplementationMethods(state.fileSet, file, normalisedPath, state.implementations)
	index.CollectImplementationBindings(
		state.fileSet,
		file,
		normalisedPath,
		&state.implementationBindings,
	)

	// Single-letter checks skip test files to reduce noise in table-driven structures and assertion
	// helpers.
	if !isTestFile {
		state.violations = append(state.violations, rules.CheckSingleLetterVars(state.fileSet, file)...)
	}
}

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	typeAwareViolations, typeAwareRan := rules.CollectTypeAwareDomainIdentifierCastViolations(
		scanRoots,
		state.scannedGoFiles,
	)
	if typeAwareRan {
		state.violations = append(state.violations, typeAwareViolations...)
	}

	state.violations = append(
		state.violations,
		index.CheckMockOrderAgainstInterfaces(state.interfaces, state.mocks)...,
	)
	state.violations = append(
		state.violations,
		index.CheckImplementationOrderAgainstInterfaces(
			state.interfaces,
			state.implementations,
			state.implementationBindings,
		)...,
	)
}

/* ------------------------------------------ Reporting ----------------------------------------- */

func sortViolations(violations []report.Violation) {
	sort.Slice(violations, func(i int, j int) bool {
		if violations[i].Position.Filename == violations[j].Position.Filename {
			return violations[i].Position.Line < violations[j].Position.Line
		}
		return violations[i].Position.Filename < violations[j].Position.Filename
	})
}

func dedupeViolations(violations []report.Violation) (deduped []report.Violation) {
	seen := make(map[string]bool)
	deduped = make([]report.Violation, 0, len(violations))

	for _, current := range violations {
		key := fmt.Sprintf(
			"%s:%d:%d|%s|%s",
			current.Position.Filename,
			current.Position.Line,
			current.Position.Column,
			current.Rule,
			current.Message,
		)

		if seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, current)
	}

	return deduped
}

func printViolations(violations []report.Violation) (exitCode int) {
	if len(violations) == 0 {
		return 0
	}

	for _, current := range violations {
		fmt.Fprintf(os.Stderr, "%s: [%s] %s\n", current.Position, current.Rule, current.Message)
	}

	return 1
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}
