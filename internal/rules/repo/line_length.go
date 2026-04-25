package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ciphera/tools/internal/styleguide"
)

const (
	lineLengthLimit    = 100
	lineLengthTabWidth = 4
)

func CheckLineLengths(
	repoRoot string,
	files []string,
) (output string, err error) {
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

			if strings.HasSuffix(path, ".sh") &&
				styleguide.HasExceptionMarker(line, styleguide.ExceptionLongLine) {
				continue
			}

			expandedLine := strings.ReplaceAll(line, "\t", strings.Repeat(" ", lineLengthTabWidth))
			if len(expandedLine) <= lineLengthLimit {
				continue
			}

			found = true
			builder.WriteString(fmt.Sprintf(
				"%s:%d (%d columns, tab width %d)\n",
				RelativePath(repoRoot, path),
				lineNumber,
				len(expandedLine),
				lineLengthTabWidth,
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
