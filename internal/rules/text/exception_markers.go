package text

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/markers"
	"ciphera/tools/internal/policy"
)

func CheckExceptionMarkers(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			marker := markers.Parse(line.Text)
			if marker.Status != markers.StatusInvalid {
				return nil
			}

			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code:    "text/exception-markers/invalid",
				File:    filewalk.RelativePath(repoRoot, path),
				Line:    line.Number,
				Message: "invalid exception marker",
			})
			return nil
		})
		if err != nil {
			return contract.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}
