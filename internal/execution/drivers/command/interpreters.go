package command

import (
	"strings"

	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/process"
	"ciphera/tools/internal/style"
)

// Exit codes that linters use to signal findings rather than failure.
const (
	// ExitFindings is the conventional Unix linter findings exit code (shellcheck, markdownlint,
	// shfmt -d).
	ExitFindings = 1
	// ExitFindingsMisspell is misspell's findings exit code when invoked with -error.
	ExitFindingsMisspell = 2
)

// InterpretPlainText returns a FileInterpreter for tools whose findings output is multi-line text
// that does not decompose cleanly per line (shellcheck, markdownlint, misspell). When the tool
// exits with code, its trimmed output becomes a single diagnostic with the given code label.
func InterpretPlainText(code int, codeLabel string) (interpreter driverkit.FileInterpreter) {
	return func(result process.CommandResult) ([]style.Diagnostic, error) {
		if result.ExitCode != code {
			return nil, nil
		}

		output := strings.TrimSpace(result.Output)
		if output == "" {
			return nil, nil
		}

		return []style.Diagnostic{{
			Code:    codeLabel,
			Message: output,
		}}, nil
	}
}

// InterpretLines returns a FileInterpreter for tools whose findings output is one finding per line
// (gofmt -l, shfmt -d). When the tool exits with code, each non-empty trimmed line becomes a
// diagnostic with the given code label.
func InterpretLines(code int, codeLabel string) (interpreter driverkit.FileInterpreter) {
	return func(result process.CommandResult) ([]style.Diagnostic, error) {
		if result.ExitCode != code {
			return nil, nil
		}

		var diagnostics []style.Diagnostic
		for _, line := range strings.Split(strings.TrimSpace(result.Output), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			diagnostics = append(diagnostics, style.Diagnostic{
				Code:    codeLabel,
				Message: line,
			})
		}

		return diagnostics, nil
	}
}
