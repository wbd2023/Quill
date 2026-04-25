package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/styleguide"
)

func CheckExceptionMarkers(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	files, err := CollectFiles(repoRoot, repository, scope, ".go", ".sh")
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
			_, _, markerFound, valid := styleguide.ParseExceptionMarker(line)
			if !markerFound || valid {
				continue
			}

			found = true
			builder.WriteString(fmt.Sprintf(
				"%s:%d invalid exception marker\n",
				RelativePath(repoRoot, path),
				lineNumber,
			))
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
