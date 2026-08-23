package text

import (
	"unicode/utf8"

	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/markers"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

const nonASCIIMarker = "allow-non-ascii"

// CheckASCII scans for non-ASCII characters in text files.
func CheckASCII(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(
		repository.ResolveScopeRoots(repoRoot, scope),
		filewalk.WalkConfig{
			ExcludedDirectories: repository.ExcludedDirectories,
			GeneratedMarker:     repository.GeneratedMarker,
		},
	)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if markers.HasMarker(line.Text, nonASCIIMarker) {
				return nil
			}

			for _, character := range line.Text {
				if character > utf8.RuneSelf-1 {
					result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
						Code:    "text/ascii/non-ascii",
						File:    filewalk.DisplayPath(repoRoot, path),
						Range:   style.Range{Start: style.Position{Line: line.Number}},
						Message: "contains non-ASCII character",
					})
					break
				}
			}

			return nil
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}
