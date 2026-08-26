package architecture

import (
	"path/filepath"
	"strings"
	"testing"
)

/* ----------------------------------- Application Boundaries ----------------------------------- */

func TestEngineDoesNotImportInboundOrPresentationAdapters(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, root)
	engineDirectory := filepath.Join(root, "internal", "engine")
	for _, file := range productionGoFiles(t, engineDirectory, false, nil) {
		for _, imported := range fileImports(t, file) {
			if hasForbiddenImport(imported, modulePath,
				[]string{"internal/cli", "internal/report"}) {
				t.Fatalf("%s imports an adapter package %q", file, imported)
			}
		}
	}
}

func TestCommandEntrypointImportsOnlyCLIInternally(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, root)
	commandDirectory := filepath.Join(root, "cmd", "quill")
	internalPrefix := modulePath + "/internal/"
	cliImport := modulePath + "/internal/cli"

	for _, file := range productionGoFiles(t, commandDirectory, false, nil) {
		for _, imported := range fileImports(t, file) {
			if strings.HasPrefix(imported, internalPrefix) && imported != cliImport {
				t.Fatalf("%s imports internal package %q instead of internal/cli", file, imported)
			}
		}
	}
}

func TestCLIImportsOnlyApplicationAndPresentationPackages(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	assertLocalImportsAreAllowed(
		t,
		filepath.Join(root, "internal", "cli"),
		moduleImportPath(t, root),
		[]string{
			"internal/engine",
			"internal/report",
			"internal/style",
		},
	)
}

func TestReportImportsOnlyApplicationAndPresentationPackages(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	assertLocalImportsAreAllowed(
		t,
		filepath.Join(root, "internal", "report"),
		moduleImportPath(t, root),
		[]string{
			"internal/coverage",
			"internal/engine",
			"internal/style",
			"internal/toolchain",
		},
	)
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func assertLocalImportsAreAllowed(
	t *testing.T,
	directory string,
	modulePath string,
	allowed []string,
) {
	t.Helper()

	localPrefix := modulePath + "/"
	for _, file := range productionGoFiles(t, directory, false, nil) {
		for _, imported := range fileImports(t, file) {
			if !strings.HasPrefix(imported, localPrefix) {
				continue
			}

			relative := strings.TrimPrefix(imported, localPrefix)
			if localImportIsAllowed(relative, allowed) {
				continue
			}

			t.Fatalf("%s imports disallowed local package %q", file, imported)
		}
	}
}

func localImportIsAllowed(imported string, allowed []string) (matches bool) {
	for _, allowedPath := range allowed {
		if imported == allowedPath || strings.HasPrefix(imported, allowedPath+"/") {
			return true
		}
	}

	return false
}
