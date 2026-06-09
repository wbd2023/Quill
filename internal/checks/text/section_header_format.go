package text

import (
	"fmt"
	"regexp"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/style"
)

func validateSectionHeader(
	repoRoot string,
	path string,
	line filewalk.Line,
	body string,
	bodyPattern *regexp.Regexp,
) (diagnostics []style.Diagnostic) {
	relativePath := filewalk.RelativePath(repoRoot, path)
	lineWidth := visualWidth(line.Text)
	if lineWidth != sectionHeaderLength {
		diagnostics = append(diagnostics, style.Diagnostic{
			Code: "text/section-headers/format",
			File: relativePath,
			Line: line.Number,
			Message: fmt.Sprintf(
				"section header must be %d columns (got %d)",
				sectionHeaderLength,
				lineWidth,
			),
		})
	}

	if !bodyPattern.MatchString(body) {
		return append(diagnostics, style.Diagnostic{
			Code:    "text/section-headers/format",
			File:    relativePath,
			Line:    line.Number,
			Message: "malformed section header body",
		})
	}

	groups := bodyPattern.FindStringSubmatch(body)
	left := len(groups[1])
	right := len(groups[3])
	if left != right && left != right+1 {
		diagnostics = append(diagnostics, style.Diagnostic{
			Code:    "text/section-headers/format",
			File:    relativePath,
			Line:    line.Number,
			Message: "section header text is not centred with left-side precedence",
		})
	}

	return diagnostics
}

func visualWidth(value string) (width int) {
	for _, character := range value {
		if character == '\t' {
			width += 4
			continue
		}

		width++
	}

	return width
}
