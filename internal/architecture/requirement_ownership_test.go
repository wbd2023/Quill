package architecture

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/testutil"
)

func TestRequirementIDsStayOutOfImplementationCode(t *testing.T) {
	requirementIDPattern := regexp.MustCompile(`\b[0-9]+\.[0-9]+\.[a-z][a-z0-9-]*\b`)
	repoRoot := testutil.RepositoryRoot(t)
	internalRoot := filepath.Join(repoRoot, "internal")

	err := filepath.WalkDir(
		internalRoot,
		func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				if entry.Name() == "testdata" {
					return filepath.SkipDir
				}

				return nil
			}

			if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			contents, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			if match := requirementIDPattern.Find(contents); len(match) > 0 {
				relativePath, _ := filepath.Rel(repoRoot, path)
				t.Fatalf(
					"%s hardcodes requirement ID %q; bind requirements in quill.toml instead",
					filepath.ToSlash(relativePath),
					string(match),
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Fatalf("walk implementation files: %v", err)
	}
}
