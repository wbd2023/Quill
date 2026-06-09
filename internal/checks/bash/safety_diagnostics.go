package bash

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/style"
)

func bashSafetyDiagnostic(
	code string,
	repoRoot string,
	path string,
	line int,
	message string,
) (diagnostic style.Diagnostic) {
	return style.Diagnostic{
		Code:    code,
		File:    filewalk.RelativePath(repoRoot, path),
		Line:    line,
		Message: message,
	}
}
