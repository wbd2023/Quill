package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	sectionHeaderLength = 100
	minLinesForHeaders  = 100
)

var genericSectionHeaderNames = map[string]struct{}{
	"check":  {},
	"checks": {},
	"misc":   {},
	"other":  {},
}

/* ------------------------------------ Section Header Rules ------------------------------------ */

func CheckSectionHeaders(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	goHeaderPattern, shellHeaderPattern, headerBodyPattern := sectionHeaderPatterns()

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
		lineCount := 0
		headerCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineCount++
			lineNumber++

			body, isHeader := extractSectionHeaderBody(
				path,
				line,
				goHeaderPattern,
				shellHeaderPattern,
			)
			if !isHeader {
				continue
			}

			headerCount++
			lineWidth := visualWidth(line)
			if lineWidth != sectionHeaderLength {
				found = true
				fmt.Fprintf(
					&builder,
					"%s:%d section header not %d columns (got %d)\n",
					RelativePath(repoRoot, path),
					lineNumber,
					sectionHeaderLength,
					lineWidth,
				)
			}

			if !headerBodyPattern.MatchString(body) {
				found = true
				fmt.Fprintf(
					&builder,
					"%s:%d malformed section header body\n",
					RelativePath(repoRoot, path),
					lineNumber,
				)
				continue
			}

			groups := headerBodyPattern.FindStringSubmatch(body)
			left := len(groups[1])
			right := len(groups[3])
			if left != right && left != right+1 {
				found = true
				fmt.Fprintf(
					&builder,
					"%s:%d section header text is not centred with left-side precedence\n",
					RelativePath(repoRoot, path),
					lineNumber,
				)
			}
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}

		if lineCount >= minLinesForHeaders && headerCount == 0 {
			found = true
			fmt.Fprintf(
				&builder,
				"%s missing section headers in %d+ line file\n",
				RelativePath(repoRoot, path),
				minLinesForHeaders,
			)
		}
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
}

func CheckSectionHeaderNames(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	goHeaderPattern, shellHeaderPattern, headerBodyPattern := sectionHeaderPatterns()

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
			body, isHeader := extractSectionHeaderBody(
				path,
				scanner.Text(),
				goHeaderPattern,
				shellHeaderPattern,
			)
			if !isHeader {
				continue
			}

			title, ok := extractSectionHeaderTitle(body, headerBodyPattern)
			if !ok {
				continue
			}

			if _, disallowed := genericSectionHeaderNames[strings.ToLower(title)]; !disallowed {
				continue
			}

			found = true
			fmt.Fprintf(
				&builder,
				"%s:%d generic section header name %q; prefer a specific heading\n",
				RelativePath(repoRoot, path),
				lineNumber,
				title,
			)
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

func visualWidth(value string) (width int) {
	for _, current := range value {
		if current == '\t' {
			width += 4
			continue
		}

		width++
	}

	return width
}

/* --------------------------------------- Header Parsing --------------------------------------- */

func sectionHeaderPatterns() (
	goHeaderPattern *regexp.Regexp,
	shellHeaderPattern *regexp.Regexp,
	headerBodyPattern *regexp.Regexp,
) {
	return regexp.MustCompile(`^/\* -+ .+ -+ \*/$`),
		regexp.MustCompile(`^# -+ .+ -+$`),
		regexp.MustCompile(`^(-+) (.+) (-+)$`)
}

func extractSectionHeaderBody(
	path string,
	line string,
	goHeaderPattern *regexp.Regexp,
	shellHeaderPattern *regexp.Regexp,
) (body string, isHeader bool) {
	switch {
	case strings.HasSuffix(path, ".go"):
		if !goHeaderPattern.MatchString(line) {
			return "", false
		}

		body = strings.TrimPrefix(line, "/* ")
		body = strings.TrimSuffix(body, " */")
		return body, true

	case strings.HasSuffix(path, ".sh"):
		if !shellHeaderPattern.MatchString(line) {
			return "", false
		}

		return strings.TrimPrefix(line, "# "), true

	default:
		return "", false
	}
}

func extractSectionHeaderTitle(
	body string,
	headerBodyPattern *regexp.Regexp,
) (title string, ok bool) {
	if !headerBodyPattern.MatchString(body) {
		return "", false
	}

	groups := headerBodyPattern.FindStringSubmatch(body)
	return groups[2], true
}
