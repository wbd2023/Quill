package architecture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDriversHaveNoLegacyFamilyImports(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, toolsRoot)
	legacyFamilies := []string{
		modulePath + "/internal/execution/drivers/command",
		modulePath + "/internal/execution/drivers/profile",
		modulePath + "/internal/execution/drivers/scan",
		modulePath + "/internal/execution/drivers/target",
		modulePath + "/internal/execution/drivers/internal",
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
				if !isDriverFamilyImport(imported, legacyFamilies) {
					continue
				}

				t.Fatalf("%s imports removed Driver package %q", path, imported)
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
