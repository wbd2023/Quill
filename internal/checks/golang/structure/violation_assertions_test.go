package structure

import (
	"strings"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func hasViolation(violations []analysis.Violation, rule string) (found bool) {
	for _, violation := range violations {
		if violation.Rule == rule {
			return true
		}
	}

	return false
}

func hasViolationAt(
	violations []analysis.Violation,
	rule string,
	line int,
	messageFragment string,
) (found bool) {
	for _, violation := range violations {
		if violation.Rule != rule || violation.Position.Line != line {
			continue
		}

		if strings.Contains(violation.Message, messageFragment) {
			return true
		}
	}

	return false
}
