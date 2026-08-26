package bash

import "strings"

func (state *shellSafetyState) checkScriptShape(
	root string,
	path string,
	patterns safetyPatterns,
	lineNumber int,
	line string,
	trimmed string,
) {
	if patterns.which.MatchString(line) {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/script-shape",
			root,
			path,
			lineNumber,
			"detect dependencies with command -v, not which",
		))
	}

	if looksLikeManualTempPath(trimmed) {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/temp-path",
			root,
			path,
			lineNumber,
			"temporary resources must be created with mktemp",
		))
	}

	if patterns.readLoop.MatchString(line) {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/script-shape",
			root,
			path,
			lineNumber,
			"avoid cmd | while read loops when loop state must survive",
		))
	}
}

func (state *shellSafetyState) addScriptShapeDiagnostics(
	root string,
	path string,
	lines []string,
) {
	if !isNonTrivialShellScript(state.functions) {
		return
	}

	lastFunction := state.functions[len(state.functions)-1]
	if lastFunction.name != "main" {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/script-shape",
			root,
			path,
			lastFunction.line,
			"non-trivial Bash scripts must keep main() as the bottom-most function",
		))
	}

	if lastLine := lastSignificantShellLine(lines); lastLine != `main "$@"` {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/script-shape",
			root,
			path,
			0,
			`non-trivial Bash scripts must end with main "$@"`,
		))
	}
}

func isNonTrivialShellScript(functions []shellFunction) (found bool) {
	for _, function := range functions {
		if function.name != "main" {
			return true
		}
	}

	return false
}

func lastSignificantShellLine(lines []string) (line string) {
	for index := len(lines) - 1; index >= 0; index-- {
		trimmed := strings.TrimSpace(lines[index])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		return trimmed
	}

	return ""
}
