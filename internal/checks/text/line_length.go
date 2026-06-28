package text

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/markers"
	"ciphera/tools/internal/style"
)

// line_length constants.
const (
	lineLengthLimit    = 100
	lineLengthTabWidth = 4
	longLineMarker     = "allow-long-line"
)

// CheckLineLengths check line lengths.
func CheckLineLengths(
	repoRoot string,
	files []string,
) (result style.ExecutionResult, err error) {
	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if strings.HasSuffix(path, ".sh") &&
				markers.HasMarker(line.Text, longLineMarker) {
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

			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
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
			return style.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}
