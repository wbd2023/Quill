package text

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/styleguide/markers"
)

const (
	lineLengthLimit    = 100
	lineLengthTabWidth = 4
)

func CheckLineLengths(
	repoRoot string,
	files []string,
) (result contract.ExecutionResult, err error) {
	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if strings.HasSuffix(path, ".sh") &&
				markers.Has(line.Text, markers.LongLine) {
				return nil
			}

			expandedLine := strings.ReplaceAll(
				line.Text,
				"\t",
				strings.Repeat(" ", lineLengthTabWidth),
			)
			if len(expandedLine) <= lineLengthLimit {
				return nil
			}

			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/line-length/too-long",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"%d columns, tab width %d",
					len(expandedLine),
					lineLengthTabWidth,
				),
			})
			return nil
		})
		if err != nil {
			return contract.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}
