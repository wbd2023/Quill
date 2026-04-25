package gostyle

import (
	"fmt"
	"path/filepath"
	"sort"

	"ciphera/tools/internal/rules/go/checks"
)

/* ------------------------------------------ Reporting ----------------------------------------- */

func sortViolations(violations []checks.Violation) {
	sort.Slice(violations, func(i int, j int) bool {
		if violations[i].Position.Filename == violations[j].Position.Filename {
			return violations[i].Position.Line < violations[j].Position.Line
		}
		return violations[i].Position.Filename < violations[j].Position.Filename
	})
}

func dedupeViolations(violations []checks.Violation) (deduped []checks.Violation) {
	seen := make(map[string]bool)
	deduped = make([]checks.Violation, 0, len(violations))

	for _, current := range violations {
		key := fmt.Sprintf(
			"%s:%d:%d|%s|%s",
			current.Position.Filename,
			current.Position.Line,
			current.Position.Column,
			current.Rule,
			current.Message,
		)

		if seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, current)
	}

	return deduped
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}
