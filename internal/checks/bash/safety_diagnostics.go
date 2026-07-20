package bash

import (
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/style"
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
