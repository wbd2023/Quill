package architecture

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEngineDoesNotImportInboundOrPresentationAdapters(t *testing.T) {
	t.Parallel()

	repositoryRoot := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, repositoryRoot)
	engineDirectory := filepath.Join(repositoryRoot, "internal", "engine")
	for _, file := range productionGoFiles(t, engineDirectory, false, nil) {
		for _, imported := range fileImports(t, file) {
			if forbiddenImport(imported, modulePath, []string{"internal/cli", "internal/report"}) {
				t.Fatalf("%s imports an adapter package %q", file, imported)
			}
		}
	}
}

func TestCommandEntrypointImportsOnlyCLIInternally(t *testing.T) {
	t.Parallel()

	repositoryRoot := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, repositoryRoot)
	commandDirectory := filepath.Join(repositoryRoot, "cmd", "quill")
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

	repositoryRoot := importBoundaryRoot(t)
	assertLocalImportsAreAllowed(
		t,
		filepath.Join(repositoryRoot, "internal", "cli"),
		moduleImportPath(t, repositoryRoot),
		[]string{
			"internal/engine",
			"internal/report",
			"internal/style",
		},
	)
}

func TestReportImportsOnlyApplicationAndPresentationPackages(t *testing.T) {
	t.Parallel()

	repositoryRoot := importBoundaryRoot(t)
	assertLocalImportsAreAllowed(
		t,
		filepath.Join(repositoryRoot, "internal", "report"),
		moduleImportPath(t, repositoryRoot),
		[]string{
			"internal/coverage",
			"internal/engine",
			"internal/style",
			"internal/toolchain",
		},
	)
}

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
