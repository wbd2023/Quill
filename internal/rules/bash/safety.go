package bash

import (
	"os"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
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

/* ---------------------------------------- Safety Rules ---------------------------------------- */

func CheckSafety(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	shellFunctionPattern := regexp.MustCompile(
		`^\s*(?:function\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*(?:\(\))?\s*\{`,
	)
	shellAssignmentPattern := regexp.MustCompile(
		`^\s*(?:local\s+)?([A-Za-z_][A-Za-z0-9_]*)=`,
	)
	shellExportPattern := regexp.MustCompile(
		`^\s*(?:readonly\s+|export\s+)([A-Za-z_][A-Za-z0-9_]*)=`,
	)
	shellWhichPattern := regexp.MustCompile(`\bwhich\s+[A-Za-z0-9_.-]+`)
	shellReadLoopPattern := regexp.MustCompile(`\|\s*while\b.*\bread\b`)
	shellcheckDisablePattern := regexp.MustCompile(
		`^\s*#\s*shellcheck\s+disable=([A-Z0-9,]+)(?:\s+--\s+(.+))?\s*$`,
	)

	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return contract.ExecutionResult{}, readErr
		}

		lines := strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n")
		functions := make([]shellFunction, 0)
		foundMktemp := false
		foundTrap := false

		for index, line := range lines {
			lineNumber := index + 1
			trimmed := strings.TrimSpace(line)

			if strings.Contains(line, "mktemp") {
				foundMktemp = true
			}
			if strings.Contains(trimmed, "trap ") {
				foundTrap = true
			}

			if matches := shellFunctionPattern.FindStringSubmatch(line); len(matches) > 1 {
				name := matches[1]
				functions = append(functions, shellFunction{line: lineNumber, name: name})

				if name != "main" && !isLowerSnakeCase(name) {
					result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
						"bash/safety/naming",
						repoRoot,
						path,
						lineNumber,
						"Bash function names should use lower-case with underscores",
					))
				}
			}

			if matches := shellExportPattern.FindStringSubmatch(line); len(matches) > 1 &&
				!isUpperSnakeCase(matches[1]) {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/naming",
					repoRoot,
					path,
					lineNumber,
					"Bash constants and exported variables should use upper-case with underscores",
				))
			}

			if matches := shellAssignmentPattern.FindStringSubmatch(line); len(matches) > 1 {
				name := matches[1]
				if !isUpperSnakeCase(name) && !isLowerSnakeCase(name) {
					result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
						"bash/safety/naming",
						repoRoot,
						path,
						lineNumber,
						"Bash non-exported variable names should use lower-case with underscores",
					))
				}
			}

			if shellWhichPattern.MatchString(line) {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/script-shape",
					repoRoot,
					path,
					lineNumber,
					"detect dependencies with command -v, not which",
				))
			}

			if looksLikeManualTempPath(trimmed) {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/temp-path",
					repoRoot,
					path,
					lineNumber,
					"temporary resources must be created with mktemp",
				))
			}

			if shellReadLoopPattern.MatchString(line) {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/script-shape",
					repoRoot,
					path,
					lineNumber,
					"avoid cmd | while read loops when loop state must survive",
				))
			}

			if strings.Contains(trimmed, "shellcheck disable=") &&
				!hasLocalShellcheckSuppressionReason(trimmed, shellcheckDisablePattern) {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/suppression",
					repoRoot,
					path,
					lineNumber,
					"shellcheck suppressions must include rule IDs and a short reason",
				))
			}
		}

		if foundMktemp && !foundTrap {
			result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
				"bash/safety/temp-path",
				repoRoot,
				path,
				0,
				"Bash scripts using mktemp must install trap-based cleanup",
			))
		}

		if isNonTrivialShellScript(functions) {
			lastFunction := functions[len(functions)-1]
			if lastFunction.name != "main" {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/script-shape",
					repoRoot,
					path,
					lastFunction.line,
					"non-trivial Bash scripts must keep main() as the bottom-most function",
				))
			}

			if lastLine := lastSignificantShellLine(lines); lastLine != `main "$@"` {
				result.Diagnostics = append(result.Diagnostics, bashSafetyDiagnostic(
					"bash/safety/script-shape",
					repoRoot,
					path,
					0,
					`non-trivial Bash scripts must end with main "$@"`,
				))
			}
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

func bashSafetyDiagnostic(
	code string,
	repoRoot string,
	path string,
	line int,
	message string,
) (diagnostic contract.Diagnostic) {
	return contract.Diagnostic{
		Code:    code,
		File:    filewalk.RelativePath(repoRoot, path),
		Line:    line,
		Message: message,
	}
}

/* ------------------------------------------- Naming ------------------------------------------- */

func isLowerSnakeCase(value string) (found bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character == '_' ||
			('a' <= character && character <= 'z') ||
			('0' <= character && character <= '9') {
			continue
		}

		return false
	}

	return true
}

func isUpperSnakeCase(value string) (found bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character == '_' ||
			('A' <= character && character <= 'Z') ||
			('0' <= character && character <= '9') {
			continue
		}

		return false
	}

	return true
}

/* --------------------------------------- Temporary Paths -------------------------------------- */

func looksLikeManualTempPath(line string) (found bool) {
	if strings.Contains(line, "mktemp") {
		return false
	}

	if strings.Contains(line, "/tmp/") || strings.Contains(line, "/var/tmp/") {
		return true
	}

	return strings.Contains(line, "TMPDIR=") || strings.Contains(line, "tmp_dir=/tmp")
}

/* ---------------------------------------- Suppressions ---------------------------------------- */

func hasLocalShellcheckSuppressionReason(
	line string,
	shellcheckDisablePattern *regexp.Regexp,
) (found bool) {
	matches := shellcheckDisablePattern.FindStringSubmatch(line)
	if len(matches) < shellcheckMatchesLength {
		return false
	}

	return matches[shellcheckRuleCaptureIndex] != "" &&
		strings.TrimSpace(matches[shellcheckReasonCaptureIndex]) != ""
}

/* ---------------------------------------- Script Shape ---------------------------------------- */

func isNonTrivialShellScript(functions []shellFunction) (found bool) {
	for _, function := range functions {
		if function.name != "main" {
			return true
		}
	}

	return false
}

func lastSignificantShellLine(lines []string) (line string) {
	for index := len(lines) - 1; index >= 0; index-- {
		trimmed := strings.TrimSpace(lines[index])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		return trimmed
	}

	return ""
}
