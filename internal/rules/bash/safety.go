package bashstyle

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	repostyle "ciphera/tools/internal/rules/repo"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	shellcheckRuleCaptureIndex   = 1
	shellcheckReasonCaptureIndex = 2
	shellcheckMatchesLength      = 3
)

/* -------------------------------------------- Types ------------------------------------------- */

type shellFunction struct {
	line int
	name string
}

/* ---------------------------------------- Safety Rules ---------------------------------------- */

func CheckSafety(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
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

	files, err := repostyle.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	foundViolation := false

	for _, path := range files {
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", readErr
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
					foundViolation = true
					builder.WriteString(fmt.Sprintf(
						"%s:%d Bash function names should use lower-case with underscores\n",
						repostyle.RelativePath(repoRoot, path),
						lineNumber,
					))
				}
			}

			if matches := shellExportPattern.FindStringSubmatch(line); len(matches) > 1 &&
				!isUpperSnakeCase(matches[1]) {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d Bash constants and exported variables should use "+
						"upper-case with underscores\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if matches := shellAssignmentPattern.FindStringSubmatch(line); len(matches) > 1 {
				name := matches[1]
				if !isUpperSnakeCase(name) && !isLowerSnakeCase(name) {
					foundViolation = true
					builder.WriteString(fmt.Sprintf(
						"%s:%d Bash non-exported variable names should use "+
							"lower-case with underscores\n",
						repostyle.RelativePath(repoRoot, path),
						lineNumber,
					))
				}
			}

			if shellWhichPattern.MatchString(line) {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d detect dependencies with command -v, not which\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if looksLikeManualTempPath(trimmed) {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d temporary resources must be created with mktemp\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if shellReadLoopPattern.MatchString(line) {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d avoid cmd | while read loops when loop state must survive\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if strings.Contains(trimmed, "shellcheck disable=") &&
				!hasLocalShellcheckSuppressionReason(trimmed, shellcheckDisablePattern) {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d shellcheck suppressions must include rule IDs and a short reason\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}
		}

		if foundMktemp && !foundTrap {
			foundViolation = true
			builder.WriteString(fmt.Sprintf(
				"%s Bash scripts using mktemp must install trap-based cleanup\n",
				repostyle.RelativePath(repoRoot, path),
			))
		}

		if isNonTrivialShellScript(functions) {
			lastFunction := functions[len(functions)-1]
			if lastFunction.name != "main" {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d non-trivial Bash scripts must keep main() as the bottom-most function\n",
					repostyle.RelativePath(repoRoot, path),
					lastFunction.line,
				))
			}

			if lastLine := lastSignificantShellLine(lines); lastLine != `main "$@"` {
				foundViolation = true
				builder.WriteString(fmt.Sprintf(
					"%s non-trivial Bash scripts must end with main \"$@\"\n",
					repostyle.RelativePath(repoRoot, path),
				))
			}
		}
	}

	if !foundViolation {
		return "", nil
	}

	return builder.String(), errViolationsFound
}

/* ------------------------------------------- Naming ------------------------------------------- */

func isLowerSnakeCase(value string) (found bool) {
	if value == "" {
		return false
	}

	for _, current := range value {
		if current == '_' ||
			('a' <= current && current <= 'z') ||
			('0' <= current && current <= '9') {
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

	for _, current := range value {
		if current == '_' ||
			('A' <= current && current <= 'Z') ||
			('0' <= current && current <= '9') {
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
