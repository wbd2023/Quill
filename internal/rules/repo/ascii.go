package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/styleguide"
)

func CheckASCII(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	files, err := CollectAllFiles(repoRoot, repository, scope)
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
			if styleguide.HasExceptionMarker(line, styleguide.ExceptionNonASCII) {
				continue
			}

			for _, runeValue := range line {
				if runeValue > utf8.RuneSelf-1 {
					found = true
					builder.WriteString(fmt.Sprintf(
						"%s:%d\n",
						RelativePath(repoRoot, path),
						lineNumber,
					))
					break
				}
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
