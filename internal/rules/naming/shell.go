package naming

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func checkShellNaming(
	result *contract.ExecutionResult,
	repoRoot string,
	path string,
	naming policy.NamingConfig,
) (err error) {
	shellAssignmentPattern := compileShellAssignmentPattern(naming.ShellForbiddenAssignments)

	return filewalk.ScanLines(path, func(line filewalk.Line) error {
		name := matchedShellAssignment(shellAssignmentPattern, line.Text)
		if name == "" {
			return nil
		}

		result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
			Code: "naming/vocabulary/shell-assignment",
			File: filewalk.RelativePath(repoRoot, path),
			Line: line.Number,
			Message: fmt.Sprintf(
				"use descriptive constant names in Bash (prefer %s over %s)",
				naming.ShellPreferredAssignment,
				name,
			),
		})
		return nil
	})
}

func compileShellAssignmentPattern(names []string) (pattern *regexp.Regexp) {
	if len(names) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`(^|[[:space:]])(local[[:space:]]+)?(%s)=`, strings.Join(names, "|")),
	)
}

func matchedShellAssignment(pattern *regexp.Regexp, line string) (name string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < shellAssignmentMatchLength {
		return ""
	}

	return matches[3]
}
