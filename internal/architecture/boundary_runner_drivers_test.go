package architecture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerDriverFamiliesStayBehindFacade(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	families := []string{
		"ciphera/tools/internal/runner/drivers/command",
		"ciphera/tools/internal/runner/drivers/project",
		"ciphera/tools/internal/runner/drivers/scan",
		"ciphera/tools/internal/runner/drivers/target",
	}

	err := filepath.WalkDir(
		filepath.Join(toolsRoot, "internal"),
		func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() || !isProductionGoFile(path) {
				return nil
			}

			relative, err := filepath.Rel(toolsRoot, filepath.Dir(path))
			if err != nil {
				return err
			}

			if filepath.ToSlash(relative) == "internal/runner/drivers" {
				return nil
			}

			for _, imported := range fileImports(t, path) {
				if !isDriverFamilyImport(imported, families) {
					continue
				}

				t.Fatalf(
					"%s imports driver family package %q; import runner/drivers facade instead",
					path,
					imported,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Fatalf("walk internal packages: %v", err)
	}
}

func isDriverFamilyImport(imported string, families []string) (found bool) {
	for _, family := range families {
		if imported == family || strings.HasPrefix(imported, family+"/") {
			return true
		}
	}

	return false
}
