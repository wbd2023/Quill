package vocabulary

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/style"
)

func checkBashVocabulary(
	result *style.ExecutionResult,
	repoRoot string,
	path string,
	config vocabularypolicy.Config,
) (err error) {
	variablePreferred := flattenSuffixMap(config.Bash.VariableNames)
	bashAssignmentPattern := compileBashAssignmentPattern(variablePreferred)

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
				"bash variable %q must be %s",
				name,
				variablePreferred[name],
			),
		})
		return nil
	})
}

func compileBashAssignmentPattern(names map[string]string) (pattern *regexp.Regexp) {
	if len(names) == 0 {
		return nil
	}

	forbidden := make([]string, 0, len(names))
	for name := range names {
		forbidden = append(forbidden, regexp.QuoteMeta(name))
	}

	return regexp.MustCompile(
		fmt.Sprintf(`(^|[[:space:]])(local[[:space:]]+)?(%s)=`, strings.Join(forbidden, "|")),
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
