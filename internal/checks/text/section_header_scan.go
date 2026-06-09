package text

import (
	"regexp"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/style"
)

func scanSectionHeaders(
	repoRoot string,
	path string,
	patterns sectionHeaderPatterns,
) (lineCount int, headers []sectionHeader, diagnostics []style.Diagnostic, err error) {
	err = filewalk.ScanLines(path, func(line filewalk.Line) error {
		lineCount++
		body, isHeader := extractSectionHeaderBody(path, line.Text, patterns.Go, patterns.Shell)
		if !isHeader {
			return nil
		}

		headers = append(headers, sectionHeader{
			Body: body,
			Line: line.Number,
		})

		diagnostics = append(
			diagnostics,
			validateSectionHeader(repoRoot, path, line, body, patterns.Body)...,
		)
		return nil
	})
	return lineCount, headers, diagnostics, err
}

func newSectionHeaderPatterns() (patterns sectionHeaderPatterns) {
	return sectionHeaderPatterns{
		Go:    regexp.MustCompile(`^/\* -+ .+ -+ \*/$`),
		Shell: regexp.MustCompile(`^# -+ .+ -+$`),
		Body:  regexp.MustCompile(`^(-+) (.+) (-+)$`),
	}
}

func extractSectionHeaderBody(
	path string,
	line string,
	goHeaderPattern *regexp.Regexp,
	shellHeaderPattern *regexp.Regexp,
) (body string, ok bool) {
	switch {
	case strings.HasSuffix(path, ".go") && goHeaderPattern.MatchString(line):
		return strings.TrimSuffix(strings.TrimPrefix(line, "/* "), " */"), true
	case strings.HasSuffix(path, ".sh") && shellHeaderPattern.MatchString(line):
		return strings.TrimPrefix(line, "# "), true
	default:
		return "", false
	}
}

func extractSectionHeaderTitle(
	body string,
	headerBodyPattern *regexp.Regexp,
) (title string, ok bool) {
	groups := headerBodyPattern.FindStringSubmatch(body)
	if len(groups) != sectionHeaderMatchGroups {
		return "", false
	}

	return strings.TrimSpace(groups[2]), true
}
