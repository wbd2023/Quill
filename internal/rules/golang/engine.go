// Package golang performs AST-based style checks on Go source files.
//
// It enforces the following rules from STYLE.md:
//   - 3.3 Named returns: all functions must use named, descriptive return values.
//   - 3.3 Naked returns: explicit return values are required.
//   - 3.3 Type elision: each parameter must have its own type.
//   - 3.3 Domain ID constructors: avoid direct casts for key domain identifier types.
//     This uses a type-aware pass with syntax fallback for non-buildable snippets.
//   - 3.1 Error handling: lowercase/no-punctuation error context, no secrets in fmt.Errorf args,
//     and sentinel errors scoped to domain/errors.go.
//   - 3.1 Adapter error wrapping: reject bare `return err` propagation in adapters.
//   - 2.2 Structured logging: slog-style calls must use stable lower-case ASCII key/value fields
//     and must not log secrets.
//   - 2.3 Sensitive data: reject hard-coded secret-like string literals in production code.
//   - 2.4 Cryptography: reject math/rand and deprecated crypto imports in production code.
//   - 2.5 Process execution: reject exec.Command shell interpolation via `sh -c` and friends.
//   - 3.2 Context/resources: reject context fields, missing HTTP client timeouts, and ignored
//     close errors without inline justification.
//   - 3.4 Inline comment style: trailing comments must start lower-case and avoid punctuation.
//   - 3.3 Single-letter variables: only i, j, k (loops) and receivers.
//   - 3.5 CRUD-L ordering inside interfaces.
//   - 3.5 Mock method order matches interface method order exactly.
//   - 3.5 Implementation method order matches interface method order exactly.
//   - 3.7 Parameter ordering: ctx first, secrets last.
//   - 3.8 Constructor ordering: profile-defined dependency categories.
//   - 3.9 File structure ordering for top-level declarations.
//   - 3.10 Data usage: prefer named struct literals, len-based slice emptiness checks, and
//     direct interface values instead of *interface forms.
//   - 6.1 Test hygiene: helpers call Helper, tests use t.Setenv and t.TempDir patterns.
package golang

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rules/golang/checks"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var errViolationsFound = errors.New("violations found")

/* -------------------------------------- Directory Checks -------------------------------------- */

// CheckDirectories runs the Go style checks for the provided directories.
func CheckDirectories(
	repoRoot string,
	directories []string,
	config policy.Config,
	checkNames ...string,
) (result contract.ExecutionResult, err error) {
	if err = validateScanRoots(directories); err != nil {
		return contract.ExecutionResult{}, err
	}

	if len(checkNames) == 1 {
		switch checkNames[0] {
		case GoCheckGuardClauseSpacing:
			return CheckGuardClauseSpacing(repoRoot, directories, config.Repository)
		case GoCheckSwitchCaseSpacing:
			return CheckSwitchCaseSpacing(repoRoot, directories, config.Repository)
		}
	}

	violations := analyseDirectories(repoRoot, directories, config, checkNames)
	if len(violations) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return contract.ExecutionResult{
		Diagnostics: diagnosticsFromViolations(repoRoot, violations),
	}, errViolationsFound
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
	config policy.Config,
	checkNames []string,
) (violations []checks.Violation) {
	state := newAnalysisState(repoRoot, config, checkNames)

	files, err := goFilesInDirectories(directories, config.Repository)
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
	violations []checks.Violation,
) (diagnostics []contract.Diagnostic) {
	diagnostics = make([]contract.Diagnostic, 0, len(violations))
	for _, violation := range violations {
		diagnostics = append(diagnostics, diagnosticFromViolation(repoRoot, violation))
	}

	return diagnostics
}

func diagnosticFromViolation(
	repoRoot string,
	violation checks.Violation,
) (diagnostic contract.Diagnostic) {
	path := violation.Position.Filename
	if repoRoot != "" && filepath.IsAbs(path) {
		path = filewalk.RelativePath(repoRoot, path)
	}

	return contract.Diagnostic{
		Code:    violation.Rule,
		File:    filepath.ToSlash(path),
		Line:    violation.Position.Line,
		Column:  violation.Position.Column,
		Message: violation.Message,
	}
}
