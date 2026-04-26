package text

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

const (
	sectionHeaderLength      = 100
	sectionHeaderMatchGroups = 4
)

type sectionHeader struct {
	Body string
	Line int
}

type sectionHeaderPatterns struct {
	Go    *regexp.Regexp
	Shell *regexp.Regexp
	Body  *regexp.Regexp
}

/* ---------------------------------------- Header Checks --------------------------------------- */

func CheckSectionHeaders(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders policy.SectionHeaderConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		lineCount, headers, diagnostics, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return contract.ExecutionResult{}, err
		}

		result.Diagnostics = append(result.Diagnostics, diagnostics...)
		if lineCount >= sectionHeaders.RequiredMinLines && len(headers) == 0 {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-headers/missing",
				File: filewalk.RelativePath(repoRoot, path),
				Message: fmt.Sprintf(
					"missing section headers in %d+ line file",
					sectionHeaders.RequiredMinLines,
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

func CheckSectionHeaderNames(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders policy.SectionHeaderConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	genericNames := sectionHeaderNameSet(sectionHeaders.GenericNames)

	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		_, headers, _, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return contract.ExecutionResult{}, err
		}

		for _, header := range headers {
			title, ok := extractSectionHeaderTitle(header.Body, patterns.Body)
			if !ok || !genericNames[strings.ToLower(title)] {
				continue
			}

			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-header-names/generic",
				File: filewalk.RelativePath(repoRoot, path),
				Line: header.Line,
				Message: fmt.Sprintf(
					"generic section header name %q; prefer a specific heading",
					title,
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

func CheckSectionHeaderDensity(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders policy.SectionHeaderConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		lineCount, headers, _, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return contract.ExecutionResult{}, err
		}

		relativePath := filewalk.RelativePath(repoRoot, path)
		if lineCount <= sectionHeaders.ShortFileMaxLines && len(headers) > 0 {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-header-density/short-file",
				File: relativePath,
				Message: fmt.Sprintf(
					"short %d-line file has section headers; remove them unless "+
						"they reduce navigation cost",
					lineCount,
				),
			})
		}

		if len(headers) >= sectionHeaders.OveruseCount {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-header-density/too-many",
				File: relativePath,
				Message: fmt.Sprintf(
					"%d section headers in one file; split the file or reduce header density",
					len(headers),
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

/* ---------------------------------------- Header Scans ---------------------------------------- */

func scanSectionHeaders(
	repoRoot string,
	path string,
	patterns sectionHeaderPatterns,
) (lineCount int, headers []sectionHeader, diagnostics []contract.Diagnostic, err error) {
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

/* ---------------------------------------- Header Format --------------------------------------- */

func validateSectionHeader(
	repoRoot string,
	path string,
	line filewalk.Line,
	body string,
	bodyPattern *regexp.Regexp,
) (diagnostics []contract.Diagnostic) {
	relativePath := filewalk.RelativePath(repoRoot, path)
	lineWidth := visualWidth(line.Text)
	if lineWidth != sectionHeaderLength {
		diagnostics = append(diagnostics, contract.Diagnostic{
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
		return append(diagnostics, contract.Diagnostic{
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
		diagnostics = append(diagnostics, contract.Diagnostic{
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

/* --------------------------------------- Header Parsing --------------------------------------- */

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

func sectionHeaderNameSet(names []string) (set map[string]bool) {
	set = make(map[string]bool, len(names))
	for _, name := range names {
		set[strings.ToLower(name)] = true
	}

	return set
}
