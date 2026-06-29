package target

import (
	"strings"

	"ciphera/tools/internal/style"
)

// appendDiagnostics converts non-empty tool output into a diagnostic and appends it. Used by
// target drivers that run external tools and need their findings as diagnostics.
func appendDiagnostics(
	diagnostics []style.Diagnostic,
	output string,
	code string,
) (result []style.Diagnostic) {
	output = strings.TrimSpace(output)
	if output == "" {
		return diagnostics
	}

	return append(diagnostics, style.Diagnostic{
		Code:    code,
		Message: output,
	})
}
