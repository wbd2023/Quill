package bash

/* ------------------------------------- Safety Diagnostics ------------------------------------- */

func (state *shellSafetyState) checkFunctionName(
	repoRoot string,
	path string,
	patterns safetyPatterns,
	lineNumber int,
	line string,
) {
	matches := patterns.function.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return
	}

	name := matches[1]
	state.functions = append(state.functions, shellFunction{line: lineNumber, name: name})
	if name == "main" || isLowerSnakeCase(name) {
		return
	}

	state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
		"bash/safety/naming",
		repoRoot,
		path,
		lineNumber,
		"Bash function names should use lower-case with underscores",
	))
}

func (state *shellSafetyState) checkVariableName(
	repoRoot string,
	path string,
	patterns safetyPatterns,
	lineNumber int,
	line string,
) {
	if matches := patterns.export.FindStringSubmatch(line); len(matches) > 1 &&
		!isUpperSnakeCase(matches[1]) {
		state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
			"bash/safety/naming",
			repoRoot,
			path,
			lineNumber,
			"Bash constants and exported variables should use upper-case with underscores",
		))
	}

	matches := patterns.assignment.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return
	}

	name := matches[1]
	if isUpperSnakeCase(name) || isLowerSnakeCase(name) {
		return
	}

	state.diagnostics = append(state.diagnostics, bashSafetyDiagnostic(
		"bash/safety/naming",
		repoRoot,
		path,
		lineNumber,
		"Bash non-exported variable names should use lower-case with underscores",
	))
}

/* --------------------------------------- Name Predicates -------------------------------------- */

func isLowerSnakeCase(value string) (found bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character == '_' ||
			('a' <= character && character <= 'z') ||
			('0' <= character && character <= '9') {
			continue
		}

		return false
	}

	return true
}

func isUpperSnakeCase(value string) (found bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character == '_' ||
			('A' <= character && character <= 'Z') ||
			('0' <= character && character <= '9') {
			continue
		}

		return false
	}

	return true
}
