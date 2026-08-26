package bash

import (
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/style"
)

func bashSafetyDiagnostic(
	code string,
	root string,
	path string,
	line int,
	message string,
) (diagnostic style.Diagnostic) {
	return style.Diagnostic{
		Code:    code,
		File:    filewalk.DisplayPath(root, path),
		Range:   style.Range{Start: style.Position{Line: line}},
		Message: message,
	}
}
