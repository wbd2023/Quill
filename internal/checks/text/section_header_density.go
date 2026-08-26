package text

import (
	"fmt"

	"github.com/wbd2023/quill/internal/filewalk"
	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

// CheckSectionHeaderDensity check section header density.
func CheckSectionHeaderDensity(
	root string,
	repository profile.RepositoryConfig,
	config textpolicy.SectionHeaderConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	files, err := filewalk.CollectFiles(
		repository.ResolveScopeRoots(root, scope),
		filewalk.WalkConfig{
			ExcludedDirectories: repository.ExcludedDirectories,
			GeneratedMarker:     repository.GeneratedMarker,
		},
		".go", ".sh",
	)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		lineCount, headers, _, err := scanSectionHeaders(root, path, patterns)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		relativePath := filewalk.DisplayPath(root, path)
		if lineCount <= config.ShortMaxLines && len(headers) > 0 {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "text/section-header-density/short-file",
				File: relativePath,
				Message: fmt.Sprintf(
					"short %d-line file has section headers; remove them unless "+
						"they reduce navigation cost",
					lineCount,
				),
			})
		}

		if len(headers) > config.MaxHeaderCount {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
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
		return style.ExecutionResult{}, nil
	}

	return result, nil
}
