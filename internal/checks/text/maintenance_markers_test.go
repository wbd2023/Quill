package text

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/style"
)

func TestCheckMaintenanceMarkersRejectsEmptyTodoText(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n// TODO:\nfunc Example() {}\n",
	)

	result, err := CheckMaintenanceMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected maintenance marker failure")
	}

	if !hasDiagnostic(
		result,
		"text/maintenance-markers/missing-action",
		"internal/example/example.go",
		3,
		"TODO/FIXME markers must include actionable text",
	) {
		t.Fatalf("expected actionable-text diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckMaintenanceMarkersAcceptsConcreteTodoText(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"README.md",
		"# Example\n\nTODO: document the relay retry flow\n",
	)

	result, err := CheckMaintenanceMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		style.Scope("all"),
	)
	if err != nil {
		t.Fatalf("expected valid maintenance marker to pass, diagnostics: %#v", result.Diagnostics)
	}
}
