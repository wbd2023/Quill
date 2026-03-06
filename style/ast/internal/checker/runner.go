// Package checker performs AST-based style checks on Go source files.
//
// It enforces the following rules from STYLE.md:
//   - 2.2 Named returns: all functions must use named, descriptive return values.
//   - 2.2 Naked returns: explicit return values are required.
//   - 2.2 Type elision: each parameter must have its own type.
//   - 2.2 Domain ID constructors: avoid direct casts for key domain identifier types.
//     This uses a type-aware pass with syntax fallback for non-buildable snippets.
//   - 2.1 Error handling: lowercase/no-punctuation error context, no secrets in fmt.Errorf args,
//     and sentinel errors scoped to domain/errors.go.
//   - 2.1 Adapter error wrapping: reject bare `return err` propagation in adapters.
//   - 2.3 Inline comment style: trailing comments must start lower-case and avoid punctuation.
//   - 2.2 Single-letter variables: only i, j, k (loops) and receivers.
//   - 2.2 Service package type naming: exported types end with Service/UseCase/Config.
//   - 2.5 CRUD-L ordering inside interfaces.
//   - 2.5 Mock method order matches interface method order exactly.
//   - 2.5 Implementation method order matches interface method order exactly.
//   - 2.7 Parameter ordering: ctx first, secrets last.
//   - 2.8 Constructor ordering: repos -> services -> adapters -> config -> secrets.
//   - 2.9 File structure ordering for top-level declarations.
package checker

import (
	"fmt"
	"os"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	minRequiredArgs = 2
	usageExitCode   = 2
)

/* --------------------------------------------- Run -------------------------------------------- */

// Run analyses the provided directories and returns a process-style exit code.
func Run(arguments []string) (exitCode int) {
	directories, ok := parseDirectories(arguments)
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: stylecheck <dir>...")
		return usageExitCode
	}

	state := newAnalysisState()

	for _, directory := range directories {
		state.walkDirectory(directory)
	}

	state.addCrossFileViolations(directories)
	state.violations = dedupeViolations(state.violations)
	sortViolations(state.violations)
	return printViolations(state.violations)
}

func parseDirectories(arguments []string) (directories []string, ok bool) {
	if len(arguments) < minRequiredArgs-1 {
		return nil, false
	}

	return arguments, true
}
