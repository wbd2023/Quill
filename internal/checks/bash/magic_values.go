package bash

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// magic_values constants.
const (
	trivialZero        = "0"
	trivialOne         = "1"
	trivialNegativeOne = "-1"
	firstCaptureIndex  = 1
	secondCaptureIndex = 2
)

/* -------------------------------------- Magic Value Rules ------------------------------------- */

// CheckMagicValues check magic values.
func CheckMagicValues(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	exitLiteralPattern := regexp.MustCompile(`^\s*exit\s+(-?\d+)\s*$`)
	comparisonPattern := regexp.MustCompile(`-(eq|ne|gt|lt|ge|le)\s+(-?\d+)`)
	headLimitPattern := regexp.MustCompile(`\bhead\s+-([0-9]+)\b`)

	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if shouldSkipShellNumericLine(line.Text) {
				return nil
			}

			if value := matchSingleLiteral(exitLiteralPattern, line.Text); value != "" &&
				isNonTrivialShellLiteral(value) {
				result.Diagnostics = append(result.Diagnostics, bashMagicDiagnostic(
					repoRoot,
					path,
					line.Number,
					fmt.Sprintf("uses non-trivial exit code literal %s", value),
				))
			}

			if strings.Contains(line.Text, "$#") {
				return nil
			}

			if value := matchCapturedLiteral(
				comparisonPattern,
				line.Text,
				secondCaptureIndex,
			); value != "" && isNonTrivialShellLiteral(value) {
				result.Diagnostics = append(result.Diagnostics, bashMagicDiagnostic(
					repoRoot,
					path,
					line.Number,
					fmt.Sprintf("uses non-trivial comparison literal %s", value),
				))
			}

			if value := matchCapturedLiteral(
				headLimitPattern,
				line.Text,
				firstCaptureIndex,
			); value != "" && value != trivialZero && value != trivialOne {
				result.Diagnostics = append(result.Diagnostics, bashMagicDiagnostic(
					repoRoot,
					path,
					line.Number,
					fmt.Sprintf("uses non-trivial head limit literal %s", value),
				))
			}

			return nil
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}

func bashMagicDiagnostic(
	repoRoot string,
	path string,
	line int,
	message string,
) (diagnostic style.Diagnostic) {
	return style.Diagnostic{
		Code:    "bash/magic-values/non-trivial",
		File:    filewalk.RelativePath(repoRoot, path),
		Line:    line,
		Message: message,
	}
}

/* ----------------------------------------- Line Skips ----------------------------------------- */

func shouldSkipShellNumericLine(line string) (skip bool) {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

/* -------------------------------------- Literal Matching -------------------------------------- */

func matchSingleLiteral(pattern *regexp.Regexp, line string) (value string) {
	matches := pattern.FindStringSubmatch(line)
	if len(matches) < secondCaptureIndex {
		return ""
	}

	return matches[1]
}

func matchCapturedLiteral(pattern *regexp.Regexp, line string, index int) (value string) {
	matches := pattern.FindStringSubmatch(line)
	if len(matches) <= index {
		return ""
	}

	return matches[index]
}

/* ----------------------------------- Literal Classification ----------------------------------- */

func isNonTrivialShellLiteral(literal string) (nonTrivial bool) {
	return literal != trivialZero &&
		literal != trivialOne &&
		literal != trivialNegativeOne
}
