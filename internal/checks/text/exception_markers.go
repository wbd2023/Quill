package text

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/markers"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func CheckExceptionMarkers(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			marker := markers.Parse(line.Text)
			if marker.Status != markers.StatusInvalid {
				return nil
			}

			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code:    "text/exception-markers/invalid",
				File:    filewalk.RelativePath(repoRoot, path),
				Line:    line.Number,
				Message: "invalid exception marker",
			})
			return nil
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, style.ViolationsFound()
}
