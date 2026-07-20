package bash

import (
	"path/filepath"
	"strings"

	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

// CheckTestHygiene check test hygiene.
func CheckTestHygiene(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(
		repository.ResolveScopeRoots(repoRoot, scope),
		walkConfig(repository),
	)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		if !isBashTestFile(repoRoot, path) {
			continue
		}

		foundMktemp := false
		foundTrap := false
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if strings.Contains(line.Text, "mktemp") {
				foundMktemp = true
			}

			if strings.Contains(line.Text, "trap ") {
				foundTrap = true
			}

			return nil
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}

		if foundMktemp && !foundTrap {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code:    "bash/test-hygiene/missing-cleanup",
				File:    filewalk.RelativePath(repoRoot, path),
				Message: "Bash tests using mktemp must install trap-based cleanup",
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}

func isBashTestFile(repoRoot string, path string) (found bool) {
	relativePath := filewalk.RelativePath(repoRoot, path)
	base := filepath.Base(relativePath)

	if strings.HasSuffix(base, "_test.sh") || strings.HasSuffix(base, ".bats") {
		return true
	}

	if !strings.HasSuffix(relativePath, ".sh") {
		return false
	}

	return strings.Contains(relativePath, "/test/") || strings.Contains(relativePath, "/tests/")
}
