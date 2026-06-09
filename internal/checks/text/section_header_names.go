package text

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func CheckSectionHeaderNames(
	repoRoot string,
	repository policy.RepositoryConfig,
	sectionHeaders SectionHeaderConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	patterns := newSectionHeaderPatterns()
	genericNames := sectionHeaderNameSet(sectionHeaders.GenericNames)

	files, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go", ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		_, headers, _, err := scanSectionHeaders(repoRoot, path, patterns)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		for _, header := range headers {
			title, ok := extractSectionHeaderTitle(header.Body, patterns.Body)
			if !ok || !genericNames[strings.ToLower(title)] {
				continue
			}

			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code: "text/section-header-names/generic",
				File: filewalk.RelativePath(repoRoot, path),
				Line: header.Line,
				Message: fmt.Sprintf(
					"generic section header name %q; prefer a specific heading",
					title,
				),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, style.ViolationsFound()
}

func sectionHeaderNameSet(names []string) (set map[string]bool) {
	set = make(map[string]bool, len(names))
	for _, name := range names {
		set[strings.ToLower(name)] = true
	}

	return set
}
