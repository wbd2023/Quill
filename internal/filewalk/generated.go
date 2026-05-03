package filewalk

import (
	"errors"
	"io"
	"os"
	"strings"

	"ciphera/tools/internal/policy"
)

const generatedHeaderLineLimit = 12

var generatedCommentPrefixes = []string{"//", "#", ";", "--"}

/* ------------------------------------- Generated Detection ------------------------------------ */

func isGeneratedFile(path string, repository policy.RepositoryConfig) (generated bool) {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	buffer := make([]byte, repository.GeneratedProbeBytes)
	count, readErr := file.Read(buffer)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		if closeErr := file.Close(); closeErr != nil {
			return false
		}
		return false
	}

	if closeErr := file.Close(); closeErr != nil {
		return false
	}

	return hasGeneratedHeader(string(buffer[:count]), repository.GeneratedMarker)
}

/* --------------------------------------- Header Matching -------------------------------------- */

func hasGeneratedHeader(contents string, marker string) (generated bool) {
	inspectedLines := 0
	inBlockComment := false

	for _, line := range strings.Split(contents, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		inspectedLines++
		if inspectedLines > generatedHeaderLineLimit {
			return false
		}

		if hasGeneratedHeaderLine(trimmed, marker, inBlockComment) {
			return true
		}

		if strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "<!--") {
			inBlockComment = true
		}

		if strings.Contains(trimmed, "*/") || strings.Contains(trimmed, "-->") {
			inBlockComment = false
		}
	}

	return false
}

func hasGeneratedHeaderLine(
	line string,
	marker string,
	inBlockComment bool,
) (generated bool) {
	comment := generatedCommentBody(line, inBlockComment)
	if comment == "" || !strings.Contains(comment, marker) {
		return false
	}

	return strings.Contains(strings.ToLower(comment), "generated")
}

/* --------------------------------------- Comment Parsing -------------------------------------- */

func generatedCommentBody(line string, inBlockComment bool) (comment string) {
	if inBlockComment && strings.HasPrefix(line, "*") {
		comment = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		comment = strings.TrimSuffix(comment, "*/")
		return strings.TrimSpace(comment)
	}

	if after, ok := strings.CutPrefix(line, "/*"); ok {
		comment = strings.TrimSpace(after)
		comment = strings.TrimSuffix(comment, "*/")
		return strings.TrimSpace(comment)
	}

	if after, ok := strings.CutPrefix(line, "<!--"); ok {
		comment = strings.TrimSpace(after)
		comment = strings.TrimSuffix(comment, "-->")
		return strings.TrimSpace(comment)
	}

	for _, prefix := range generatedCommentPrefixes {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		comment = strings.TrimSpace(strings.TrimPrefix(line, prefix))
		return strings.TrimSpace(comment)
	}

	return ""
}
