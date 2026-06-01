package text

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func CheckSectionHeaderDensity(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders SectionHeaderConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		lineCount, headers, _, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return contract.ExecutionResult{}, err
		}

		relativePath := filewalk.RelativePath(repoRoot, path)
		if lineCount <= sectionHeaders.ShortMaxLines && len(headers) > 0 {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-header-density/short-file",
				File: relativePath,
				Message: fmt.Sprintf(
					"short %d-line file has section headers; remove them unless "+
						"they reduce navigation cost",
					lineCount,
				),
			})
		}

		if len(headers) > sectionHeaders.MaxHeaderCount {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "text/section-header-density/too-many",
				File: relativePath,
				Message: fmt.Sprintf(
					"%d section headers in one file; split the file or reduce header density",
					len(headers),
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}
