package repostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckMaintenanceMarkersRejectsEmptyTodoText(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n// TODO:\nfunc Example() {}\n",
	)

	output, err := CheckMaintenanceMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected maintenance marker failure")
	}

	if !strings.Contains(output, "TODO/FIXME markers must include actionable text") {
		t.Fatalf("expected actionable-text violation, got:\n%s", output)
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

	output, err := CheckMaintenanceMarkers(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeAll,
	)
	if err != nil {
		t.Fatalf("expected valid maintenance marker to pass, output:\n%s", output)
	}
}
