package bash

import "strings"

func (state *shellSafetyState) addCleanupDiagnostics(repoRoot string, path string) {
	if !state.foundMktemp || state.foundTrap {
		return
	}

	state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
		"bash/safety/temp-path",
		repoRoot,
		path,
		0,
		"Bash scripts using mktemp must install trap-based cleanup",
	))
}

func looksLikeManualTempPath(line string) (found bool) {
	if strings.Contains(line, "mktemp") {
		return false
	}

	if strings.Contains(line, "/tmp/") || strings.Contains(line, "/var/tmp/") {
		return true
	}

	return strings.Contains(line, "TMPDIR=") || strings.Contains(line, "tmp_dir=/tmp")
}
