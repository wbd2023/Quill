package bash

import (
	"regexp"
	"strings"
)

func (state *shellSafetyState) checkShellcheckSuppression(
	repoRoot string,
	path string,
	patterns safetyPatterns,
	lineNumber int,
	trimmed string,
) {
	if !strings.Contains(trimmed, "shellcheck disable=") ||
		hasLocalShellcheckSuppressionReason(trimmed, patterns.shellcheckSuppression) {
		return
	}

	state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
		"bash/safety/suppression",
		repoRoot,
		path,
		lineNumber,
		"shellcheck suppressions must include rule IDs and a short reason",
	))
}

func hasLocalShellcheckSuppressionReason(
	line string,
	shellcheckDisablePattern *regexp.Regexp,
) (found bool) {
	matches := shellcheckDisablePattern.FindStringSubmatch(line)
	if len(matches) < shellcheckMatchesLength {
		return false
	}

	return matches[shellcheckRuleCaptureIndex] != "" &&
		strings.TrimSpace(matches[shellcheckReasonCaptureIndex]) != ""
}
