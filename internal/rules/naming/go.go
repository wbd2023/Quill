package naming

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

/* ------------------------------------------ Go Names ------------------------------------------ */

func checkGoNaming(
	result *contract.ExecutionResult,
	repoRoot string,
	path string,
	naming policy.NamingConfig,
) (err error) {
	goTypePattern := compileGoTypeSuffixPattern(naming.GoTypeSuffixForbidden)
	goIdentifierPattern := compileGoIdentifierSuffixPattern(naming.GoIdentifierSuffixForbidden)

	return filewalk.ScanLines(path, func(line filewalk.Line) error {
		if suffix := matchedGoTypeSuffix(goTypePattern, line.Text); suffix != "" {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "naming/vocabulary/go-type-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"use %s not %s in type names",
					naming.GoTypeSuffixPreferred,
					suffix,
				),
			})
		}

		if naming.GoIdentifierSuffixPreferred != "" &&
			strings.Contains(line.Text, naming.GoIdentifierSuffixPreferred) {
			return nil
		}

		if strings.HasPrefix(strings.TrimSpace(line.Text), "//") {
			return nil
		}

		if suffix := matchedGoIdentifierSuffix(goIdentifierPattern, line.Text); suffix != "" {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "naming/vocabulary/go-identifier-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"use x%s not x%s",
					naming.GoIdentifierSuffixPreferred,
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
