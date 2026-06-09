package vocabulary

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/style"
)

func checkBashVocabulary(
	result *style.ExecutionResult,
	repoRoot string,
	path string,
	config Config,
) (err error) {
	bashAssignmentPattern := compileBashAssignmentPattern(
		config.Bash.ForbiddenVariableNames,
	)

	return filewalk.ScanLines(path, func(line filewalk.Line) error {
		name := matchedBashAssignment(bashAssignmentPattern, line.Text)
		if name == "" {
			return nil
		}

		result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
			Code: "vocabulary/project-terms/bash-assignment",
			File: filewalk.RelativePath(repoRoot, path),
			Line: line.Number,
			Message: fmt.Sprintf(
				"use descriptive constant names in Bash (prefer %s over %s)",
				config.Bash.PreferredVariableName,
				name,
			),
		})
		return nil
	})
}

func compileBashAssignmentPattern(names []string) (pattern *regexp.Regexp) {
	if len(names) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`(^|[[:space:]])(local[[:space:]]+)?(%s)=`, strings.Join(names, "|")),
	)
}

func matchedBashAssignment(pattern *regexp.Regexp, line string) (name string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < bashAssignmentMatchLength {
		return ""
	}

	return matches[3]
}
