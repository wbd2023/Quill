package bashstyle

import (
	"bufio"
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
	trivialZeroValue        = "0"
	trivialOneValue         = "1"
	trivialNegativeOneValue = "-1"
	firstCaptureIndex       = 1
	secondCaptureIndex      = 2
)

/* -------------------------------------- Magic Value Rules ------------------------------------- */

func CheckMagicValues(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	exitLiteralPattern := regexp.MustCompile(`^\s*exit\s+(-?\d+)\s*$`)
	comparisonPattern := regexp.MustCompile(`-(eq|ne|gt|lt|ge|le)\s+(-?\d+)`)
	headLimitPattern := regexp.MustCompile(`\bhead\s+-([0-9]+)\b`)

	files, err := repostyle.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	found := false

	for _, path := range files {
		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if shouldSkipShellNumericLine(line) {
				continue
			}

			if value := matchSingleLiteral(exitLiteralPattern, line); value != "" &&
				isNonTrivialShellValue(value) {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d uses non-trivial exit code literal %s\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
					value,
				))
			}

			if strings.Contains(line, "$#") {
				continue
			}

			if value := matchCapturedLiteral(
				comparisonPattern,
				line,
				secondCaptureIndex,
			); value != "" && isNonTrivialShellValue(value) {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d uses non-trivial comparison literal %s\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
					value,
				))
			}

			if value := matchCapturedLiteral(
				headLimitPattern,
				line,
				firstCaptureIndex,
			); value != "" && value != trivialZeroValue && value != trivialOneValue {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d uses non-trivial head limit literal %s\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
					value,
				))
			}
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
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

/* ------------------------------------ Value Classification ------------------------------------ */

func isNonTrivialShellValue(value string) (nonTrivial bool) {
	return value != trivialZeroValue &&
		value != trivialOneValue &&
		value != trivialNegativeOneValue
}
