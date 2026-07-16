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
		"ciphera/tools/internal/execution/drivers/command",
		"ciphera/tools/internal/execution/drivers/profile",
		"ciphera/tools/internal/execution/drivers/scan",
		"ciphera/tools/internal/execution/drivers/target",
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

			if filepath.ToSlash(relative) == "internal/execution/drivers" {
				return nil
			}

			for _, imported := range fileImports(t, path) {
				if !isDriverFamilyImport(imported, families) {
					continue
				}

				t.Fatalf(
					"%s imports driver family package %q; import execution/drivers facade instead",
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
