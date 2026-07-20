// Package golang runs Go style checks for target directories.
package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
	"github.com/wbd2023/Quill/internal/checks/gopolicy"
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* ------------------------------------------- Errors ------------------------------------------- */

/* -------------------------------------- Directory Checks -------------------------------------- */

// CheckDirectories runs the Go style checks for the provided directories.
func CheckDirectories(
	repoRoot string,
	directories []string,
	repository policy.RepositoryConfig,
	paths policy.PathRoles,
	goConfig gopolicy.Config,
	checkNames ...string,
) (result style.ExecutionResult, err error) {
	if err = validateScanRoots(directories); err != nil {
		return style.ExecutionResult{}, err
	}

	violations := analyseDirectories(repoRoot, directories, repository, paths, goConfig, checkNames)
	if len(violations) == 0 {
		return style.ExecutionResult{}, nil
	}

	return style.ExecutionResult{
		Diagnostics: diagnosticsFromViolations(repoRoot, violations),
	}, nil
}

func validateScanRoots(directories []string) (err error) {
	for _, directory := range directories {
		info, statErr := os.Stat(directory)
		if statErr != nil {
			return fmt.Errorf("scan root %q: %w", directory, statErr)
		}

		if !info.IsDir() {
			return fmt.Errorf("scan root %q is not a directory", directory)
		}
	}

	return nil
}

func analyseDirectories(
	repoRoot string,
	directories []string,
	repository policy.RepositoryConfig,
	paths policy.PathRoles,
	goConfig gopolicy.Config,
	checkNames []string,
) (violations []analysis.Violation) {
	state := newAnalysisState(repoRoot, repository, paths, goConfig, checkNames)

	files, err := goFilesInDirectories(directories, repository)
	if err != nil {
		state.writeWarning("error walking Go files: %v\n", err)
		return nil
	}

	for _, path := range files {
		state.processFile(path)
	}

	state.addCrossFileViolations(directories)
	state.violations = dedupeViolations(state.violations)
	sortViolations(state.violations)
	return state.violations
}

func diagnosticsFromViolations(
	repoRoot string,
	violations []analysis.Violation,
) (diagnostics []style.Diagnostic) {
	diagnostics = make([]style.Diagnostic, 0, len(violations))
	for _, violation := range violations {
		diagnostics = append(diagnostics, diagnosticFromViolation(repoRoot, violation))
	}

	return diagnostics
}

func diagnosticFromViolation(
	repoRoot string,
	violation analysis.Violation,
) (diagnostic style.Diagnostic) {
	path := violation.Position.Filename
	if repoRoot != "" && filepath.IsAbs(path) {
		path = filewalk.RelativePath(repoRoot, path)
	}

	return style.Diagnostic{
		Code:    violation.Rule,
		File:    filepath.ToSlash(path),
		Line:    violation.Position.Line,
		Column:  violation.Position.Column,
		Message: violation.Message,
	}
}
