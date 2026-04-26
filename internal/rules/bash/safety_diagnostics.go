package bash

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
)

func bashSafetyDiagnostic(
	code string,
	repoRoot string,
	path string,
	line int,
	message string,
) (diagnostic contract.Diagnostic) {
	return contract.Diagnostic{
		Code:    code,
		File:    filewalk.RelativePath(repoRoot, path),
		Line:    line,
		Message: message,
	}
}
