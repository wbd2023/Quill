package text

import (
	"unicode/utf8"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/markers"
	"ciphera/tools/internal/policy"
)

const nonASCIIMarker = "allow-non-ascii"

func CheckASCII(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if markers.Has(line.Text, nonASCIIMarker) {
				return nil
			}

			for _, character := range line.Text {
				if character > utf8.RuneSelf-1 {
					result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
						Code:    "text/ascii/non-ascii",
						File:    filewalk.RelativePath(repoRoot, path),
						Line:    line.Number,
						Message: "contains non-ASCII character",
					})
					break
				}
			}

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
