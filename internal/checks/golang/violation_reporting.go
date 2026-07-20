package golang

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func sortViolations(violations []analysis.Violation) {
	sort.Slice(violations, func(i int, j int) bool {
		if violations[i].Position.Filename == violations[j].Position.Filename {
			return violations[i].Position.Line < violations[j].Position.Line
		}
		return violations[i].Position.Filename < violations[j].Position.Filename
	})
}

func dedupeViolations(violations []analysis.Violation) (deduped []analysis.Violation) {
	seen := make(map[string]bool)
	deduped = make([]analysis.Violation, 0, len(violations))

	for _, violation := range violations {
		key := fmt.Sprintf(
			"%s:%d:%d|%s|%s",
			violation.Position.Filename,
			violation.Position.Line,
			violation.Position.Column,
			violation.Rule,
			violation.Message,
		)

		if seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, violation)
	}

	return deduped
}

func normalisePath(path string) (normalisedPath string) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolute))
}
