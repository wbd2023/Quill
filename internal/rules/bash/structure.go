package bashstyle

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	repostyle "ciphera/tools/internal/rules/repo"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	shellShebangLine = "#!/bin/bash"
	strictModeLine   = "set -euo pipefail"
)

/* --------------------------------------- Structure Rules -------------------------------------- */

func CheckStructure(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
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
		hasStrictMode := false
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			if lineNumber == 1 && line != shellShebangLine {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:1 missing %s\n",
					repostyle.RelativePath(repoRoot, path),
					shellShebangLine,
				))
			}

			if line == strictModeLine {
				hasStrictMode = true
			}

			if strings.Contains(line, "\r") {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d contains CRLF line endings\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if hasTrailingHorizontalWhitespace(line) {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d has trailing whitespace\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}

			if hasSpaceIndentation(line) {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d uses space indentation\n",
					repostyle.RelativePath(repoRoot, path),
					lineNumber,
				))
			}
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}

		if !hasStrictMode {
			found = true
			builder.WriteString(fmt.Sprintf(
				"%s missing %s\n",
				repostyle.RelativePath(repoRoot, path),
				strictModeLine,
			))
		}
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
}

func hasSpaceIndentation(line string) (found bool) {
	hasLeadingSpaces := false

	for _, runeValue := range line {
		switch runeValue {
		case ' ':
			hasLeadingSpaces = true
		case '\t':
			return false
		default:
			return hasLeadingSpaces
		}
	}

	return false
}

func hasTrailingHorizontalWhitespace(line string) (found bool) {
	if line == "" {
		return false
	}

	lastByte := line[len(line)-1]
	return lastByte == ' ' || lastByte == '\t'
}
