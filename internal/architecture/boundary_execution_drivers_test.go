package architecture

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestOutsideDriverFacadeDoesNotImportDriverChildren(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, root)
	driverDirectory := "internal/execution/drivers"
	driverChildPrefix := modulePath + "/" + driverDirectory + "/"

	for _, file := range productionGoFiles(t, root, true, nil) {
		directory, err := filepath.Rel(root, filepath.Dir(file))
		if err != nil {
			t.Fatalf("rel %s: %v", file, err)
		}
		directory = filepath.ToSlash(directory)
		if directory == driverDirectory ||
			strings.HasPrefix(directory, driverDirectory+"/") {
			continue
		}

		for _, imported := range fileImports(t, file) {
			if !strings.HasPrefix(imported, driverChildPrefix) {
				continue
			}

			t.Fatalf("%s imports Driver implementation %q", file, imported)
		}
	}
}
