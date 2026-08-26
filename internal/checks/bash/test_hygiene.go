package bash

import (
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

// CheckTestHygiene check test hygiene.
func CheckTestHygiene(
	root string,
	repository profile.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(
		repository.ResolveScopeRoots(root, scope),
		walkConfig(repository),
	)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		if !isBashTestFile(root, path) {
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
				File:    filewalk.DisplayPath(root, path),
				Message: "Bash tests using mktemp must install trap-based cleanup",
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}

func isBashTestFile(root string, path string) (found bool) {
	relativePath := filewalk.DisplayPath(root, path)
	base := filepath.Base(relativePath)

	if strings.HasSuffix(base, "_test.sh") || strings.HasSuffix(base, ".bats") {
		return true
	}

	if !strings.HasSuffix(relativePath, ".sh") {
		return false
	}

	return strings.Contains(relativePath, "/test/") || strings.Contains(relativePath, "/tests/")
}
