package text

import (
	"fmt"
	"regexp"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

const (
	sectionHeaderLength      = 100
	sectionHeaderMatchGroups = 4
)

type sectionHeader struct {
	Body string
	Line int
}

type sectionHeaderPatterns struct {
	Go    *regexp.Regexp
	Shell *regexp.Regexp
	Body  *regexp.Regexp
}

func CheckSectionHeaders(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders SectionHeaderConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		lineCount, headers, diagnostics, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		result.Diagnostics = append(result.Diagnostics, diagnostics...)
		if lineCount >= sectionHeaders.LargeMinLines && len(headers) == 0 {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "text/section-headers/missing",
				File: filewalk.RelativePath(repoRoot, path),
				Message: fmt.Sprintf(
					"missing section headers in %d+ line file",
					sectionHeaders.LargeMinLines,
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, style.ViolationsFound()
}
