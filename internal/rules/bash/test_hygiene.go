package bash

import (
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func CheckTestHygiene(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return contract.ExecutionResult{}, err
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
			return contract.ExecutionResult{}, err
		}

		if foundMktemp && !foundTrap {
			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code:    "bash/test-hygiene/missing-cleanup",
				File:    filewalk.RelativePath(repoRoot, path),
				Message: "Bash tests using mktemp must install trap-based cleanup",
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
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
