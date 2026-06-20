package bash

import (
	"os"
	"regexp"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

const (
	shellcheckRuleCaptureIndex   = 1
	shellcheckReasonCaptureIndex = 2
	shellcheckMatchesLength      = 3
)

type shellFunction struct {
	line int
	name string
}

type safetyPatterns struct {
	function              *regexp.Regexp
	assignment            *regexp.Regexp
	export                *regexp.Regexp
	which                 *regexp.Regexp
	readLoop              *regexp.Regexp
	shellcheckSuppression *regexp.Regexp
}

type shellSafetyState struct {
	functions   []shellFunction
	diagnostics []style.Diagnostic
	foundMktemp bool
	foundTrap   bool
}

/* ---------------------------------------- Safety Rules ---------------------------------------- */

// CheckSafety check safety.
func CheckSafety(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	patterns := newSafetyPatterns()
	for _, path := range files {
		diagnostics, err := checkShellSafetyFile(repoRoot, path, patterns)
		if err != nil {
			return style.ExecutionResult{}, err
		}
		result.Diagnostics = append(result.Diagnostics, diagnostics...)
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, style.ViolationsFound()
}

func newSafetyPatterns() (patterns safetyPatterns) {
	return safetyPatterns{
		function: regexp.MustCompile(
			`^\s*(?:function\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*(?:\(\))?\s*\{`,
		),
		assignment: regexp.MustCompile(
			`^\s*(?:local\s+)?([A-Za-z_][A-Za-z0-9_]*)=`,
		),
		export:   regexp.MustCompile(`^\s*(?:readonly\s+|export\s+)([A-Za-z_][A-Za-z0-9_]*)=`),
		which:    regexp.MustCompile(`\bwhich\s+[A-Za-z0-9_.-]+`),
		readLoop: regexp.MustCompile(`\|\s*while\b.*\bread\b`),
		shellcheckSuppression: regexp.MustCompile(
			`^\s*#\s*shellcheck\s+disable=([A-Z0-9,]+)(?:\s+--\s+(.+))?\s*$`,
		),
	}
}

func checkShellSafetyFile(
	repoRoot string,
	path string,
	patterns safetyPatterns,
) (diagnostics []style.Diagnostic, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n")
	state := shellSafetyState{
		functions: make([]shellFunction, 0),
	}

	for index, line := range lines {
		scanShellSafetyLine(repoRoot, path, patterns, index+1, line, &state)
	}

	state.addCleanupDiagnostics(repoRoot, path)
	state.addScriptShapeDiagnostics(repoRoot, path, lines)
	return state.diagnostics, nil
}

func scanShellSafetyLine(
	repoRoot string,
	path string,
	patterns safetyPatterns,
	lineNumber int,
	line string,
	state *shellSafetyState,
) {
	trimmed := strings.TrimSpace(line)

	if strings.Contains(line, "mktemp") {
		state.foundMktemp = true
	}
	if strings.Contains(trimmed, "trap ") {
		state.foundTrap = true
	}

	state.checkFunctionName(repoRoot, path, patterns, lineNumber, line)
	state.checkVariableName(repoRoot, path, patterns, lineNumber, line)
	state.checkScriptShape(repoRoot, path, patterns, lineNumber, line, trimmed)
	state.checkShellcheckSuppression(repoRoot, path, patterns, lineNumber, trimmed)
}
