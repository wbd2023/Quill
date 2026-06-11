package vocabulary

import (
	"fmt"
	"regexp"
	"strings"

	vocabularypolicy "ciphera/tools/internal/checks/vocabulary/policy"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/style"
)

/* ------------------------------------------ Go Names ------------------------------------------ */

func checkGoVocabulary(
	result *style.ExecutionResult,
	repoRoot string,
	path string,
	config vocabularypolicy.Config,
) (err error) {
	goTypePattern := compileGoTypeSuffixPattern(config.Go.ForbiddenTypeSuffixes)
	goIdentifierPattern := compileGoIdentifierSuffixPattern(
		config.Go.ForbiddenIdentifierSuffixes,
	)

	return filewalk.ScanLines(path, func(line filewalk.Line) error {
		if suffix := matchedGoTypeSuffix(goTypePattern, line.Text); suffix != "" {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "vocabulary/project-terms/go-type-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"use %s not %s in type names",
					config.Go.PreferredTypeSuffix,
					suffix,
				),
			})
		}

		if config.Go.PreferredIdentifierSuffix != "" &&
			strings.Contains(line.Text, config.Go.PreferredIdentifierSuffix) {
			return nil
		}

		if strings.HasPrefix(strings.TrimSpace(line.Text), "//") {
			return nil
		}

		if suffix := matchedGoIdentifierSuffix(goIdentifierPattern, line.Text); suffix != "" {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "vocabulary/project-terms/go-identifier-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"use x%s not x%s",
					config.Go.PreferredIdentifierSuffix,
					suffix,
				),
			})
		}

		return nil
	})
}

/* ------------------------------------------ Patterns ------------------------------------------ */

func compileGoTypeSuffixPattern(suffixes []string) (pattern *regexp.Regexp) {
	if len(suffixes) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`type\s+\w*(%s)\s+`, strings.Join(suffixes, "|")),
	)
}

func compileGoIdentifierSuffixPattern(suffixes []string) (pattern *regexp.Regexp) {
	if len(suffixes) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`\b\w+(%s)\b`, strings.Join(suffixes, "|")),
	)
}

func matchedGoTypeSuffix(pattern *regexp.Regexp, line string) (suffix string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < goTypeSuffixMatchLength {
		return ""
	}

	return matches[1]
}

func matchedGoIdentifierSuffix(pattern *regexp.Regexp, line string) (suffix string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < goIdentifierSuffixMatchLength {
		return ""
	}

	return matches[1]
}
