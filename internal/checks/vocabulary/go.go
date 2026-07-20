package vocabulary

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/style"
)

/* ------------------------------------------ Go Names ------------------------------------------ */

func checkGoVocabulary(
	result *style.ExecutionResult,
	repoRoot string,
	path string,
	config vocabularypolicy.Config,
) (err error) {
	typePreferred := flattenSuffixMap(config.Go.TypeSuffixes)
	identifierPreferred := flattenSuffixMap(config.Go.IdentifierSuffixes)

	goTypePattern := compileGoSuffixPattern(`type\s+\w*(%s)\s+`, typePreferred)
	goIdentifierPattern := compileGoSuffixPattern(`\b\w+(%s)\b`, identifierPreferred)

	return filewalk.ScanLines(path, func(line filewalk.Line) error {
		typeSuffix := matchedSuffix(goTypePattern, line.Text, goTypeSuffixMatchLength)
		if typeSuffix != "" {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "vocabulary/project-terms/go-type-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"type name suffix %q must be %s",
					typeSuffix,
					typePreferred[typeSuffix],
				),
			})
		}

		if containsAny(line.Text, preferredKeys(config.Go.IdentifierSuffixes)) {
			return nil
		}

		if strings.HasPrefix(strings.TrimSpace(line.Text), "//") {
			return nil
		}

		idSuffix := matchedSuffix(goIdentifierPattern, line.Text, goIdentifierSuffixMatchLength)
		if idSuffix != "" {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "vocabulary/project-terms/go-identifier-suffix",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"identifier suffix %q must be %s",
					idSuffix,
					identifierPreferred[idSuffix],
				),
			})
		}

		return nil
	})
}

/* ------------------------------------------ Patterns ------------------------------------------ */

// flattenSuffixMap inverts a preferred -> [shorthands] map into a shorthand -> preferred lookup so
// a matched forbidden suffix can be resolved to its preferred form.
func flattenSuffixMap(suffixes map[string][]string) (lookup map[string]string) {
	lookup = make(map[string]string, len(suffixes))
	for preferred, forbidden := range suffixes {
		for _, shorthand := range forbidden {
			lookup[shorthand] = preferred
		}
	}

	return lookup
}

// preferredKeys returns the preferred forms (map keys) for the line-skip check. A line that already
// uses a preferred form such as Repository is not flagged for its shorthands such as Repo.
func preferredKeys(suffixes map[string][]string) (keys []string) {
	keys = make([]string, 0, len(suffixes))
	for preferred := range suffixes {
		keys = append(keys, preferred)
	}

	return keys
}

func containsAny(haystack string, needles []string) (found bool) {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}

	return false
}

func compileGoSuffixPattern(template string, suffixes map[string]string) (pattern *regexp.Regexp) {
	if len(suffixes) == 0 {
		return nil
	}

	forbidden := make([]string, 0, len(suffixes))
	for shorthand := range suffixes {
		forbidden = append(forbidden, regexp.QuoteMeta(shorthand))
	}

	return regexp.MustCompile(fmt.Sprintf(template, strings.Join(forbidden, "|")))
}

func matchedSuffix(pattern *regexp.Regexp, line string, matchLength int) (suffix string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < matchLength {
		return ""
	}

	return matches[1]
}
