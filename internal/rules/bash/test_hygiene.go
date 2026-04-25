package bashstyle

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	repostyle "ciphera/tools/internal/rules/repo"
)

/* ----------------------------------------- Bash Tests ----------------------------------------- */

func CheckTestHygiene(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	files, err := repostyle.CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	foundViolation := false

	for _, path := range files {
		if !isBashTestFile(repoRoot, path) {
			continue
		}

		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		foundMktemp := false
		foundTrap := false
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "mktemp") {
				foundMktemp = true
			}

			if strings.Contains(line, "trap ") {
				foundTrap = true
			}
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}

		if foundMktemp && !foundTrap {
			foundViolation = true
			builder.WriteString(fmt.Sprintf(
				"%s Bash tests using mktemp must install trap-based cleanup\n",
				repostyle.RelativePath(repoRoot, path),
			))
		}
	}

	if !foundViolation {
		return "", nil
	}

	return builder.String(), errViolationsFound
}

/* --------------------------------------- Test Discovery --------------------------------------- */

func isBashTestFile(repoRoot string, path string) (found bool) {
	relativePath := repostyle.RelativePath(repoRoot, path)
	base := filepath.Base(relativePath)

	if strings.HasSuffix(base, "_test.sh") || strings.HasSuffix(base, ".bats") {
		return true
	}

	if !strings.HasSuffix(relativePath, ".sh") {
		return false
	}

	return strings.Contains(relativePath, "/test/") || strings.Contains(relativePath, "/tests/")
}
