package engine

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

/* --------------------------------------- Acceptance Test -------------------------------------- */

// TestExternalPackCheckIsTheAcceptancePath proves the full MVP flow end to end: a temporary
// repository declares a local external Pack, the Profile enables and binds its rules, and a check
// runs the Pack subprocess and surfaces its structured diagnostic through the normal result model.
// It also proves the safety invariant: loading the Pack as declarations during metadata never
// launches the subprocess.
func TestExternalPackCheckIsTheAcceptancePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	stageExternalPack(t, root)

	config := policy.Profile{
		SchemaVersion: policy.SchemaVersion,
		Repository:    profiles.RepositoryConfig(),
		StyleGuide:    policy.StyleGuideConfig{Path: "STYLE.md"},
		EnabledPacks:  []string{"extpack"},
		PackSources:   []policy.PackSource{{Path: ".quill/packs/extpack"}},
		Rules: []policy.RuleBinding{
			{
				RuleID:         "extpack/forbidden-import",
				Enforcement:    style.EnforcementRequired,
				Scope:          "all",
				RequirementIDs: []string{"9.1.external-rules"},
			},
			{
				RuleID:         "extpack/marker",
				Enforcement:    style.EnforcementRequired,
				Scope:          "all",
				RequirementIDs: []string{"9.1.external-rules"},
			},
		},
	}

	testutil.WriteFile(t, root, "STYLE.md",
		"# Style Guide\n\n### 9.1 External Packs\n\n"+
			"<!-- style: id=9.1.external-rules -->\n"+
			"* External Pack rules MUST be bound to a requirement.\n")
	testutil.WriteFile(t, root, "quill.toml", profiles.Format(t, config))

	engine, err := New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	marker := filepath.Join(root, ".pack-ran-marker")

	t.Run("metadata enumerates the external pack without running it", func(t *testing.T) {
		snapshot, err := engine.Metadata(context.Background())
		if err != nil {
			t.Fatalf("Metadata: %v", err)
		}

		var seen bool
		for _, packMeta := range snapshot.Packs {
			if packMeta.ID == "extpack" {
				seen = true
				if !packMeta.External {
					t.Fatal("expected external pack to carry External provenance")
				}
			}
		}
		if !seen {
			t.Fatal("expected external pack in metadata")
		}

		if _, err := os.Stat(marker); err == nil {
			t.Fatal("external pack subprocess ran during metadata; the marker file exists")
		}
	})

	t.Run("check runs the external pack and surfaces its diagnostic", func(t *testing.T) {
		result, err := engine.Check(context.Background(), CheckOptions{})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}

		if _, err := os.Stat(marker); err != nil {
			t.Fatal("expected external pack subprocess to run during check (marker file missing)")
		}

		var diagnostic string
		var found bool
		for _, entry := range result.Rules {
			if entry.Rule.ID != "extpack/forbidden-import" {
				continue
			}
			found = true
			for _, diag := range entry.Execution.Diagnostics {
				if diag.File != "" {
					diagnostic = diag.File
				}
			}
		}
		if !found {
			t.Fatal("expected the external rule in the check result")
		}
		if diagnostic != "internal/service/users.go" {
			t.Fatalf("expected the external diagnostic file, got %q", diagnostic)
		}
	})
}

// TestExternalPackConfigReachesRuntime proves a nonempty [packs.extpack] config block is accepted
// by the external Pack's raw-config policy and forwarded unchanged to the subprocess through
// Request.Configuration, so an external Pack validates its own configuration.
func TestExternalPackConfigReachesRuntime(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	stageExternalPack(t, root)

	config := policy.Profile{
		SchemaVersion: policy.SchemaVersion,
		Repository:    profiles.RepositoryConfig(),
		StyleGuide:    policy.StyleGuideConfig{Path: "STYLE.md"},
		EnabledPacks:  []string{"extpack"},
		PackSources:   []policy.PackSource{{Path: ".quill/packs/extpack"}},
		PackConfigs: policy.PackConfigs{
			"extpack": policy.PackConfig{"allowed_packages": []any{"database/sql"}},
		},
		Rules: []policy.RuleBinding{
			{
				RuleID:         "extpack/inspect",
				Enforcement:    style.EnforcementRequired,
				Scope:          "all",
				RequirementIDs: []string{"9.1.external-rules"},
			},
		},
	}

	testutil.WriteFile(t, root, "STYLE.md",
		"# Style Guide\n\n### 9.1 External Packs\n\n"+
			"<!-- style: id=9.1.external-rules -->\n"+
			"* External Pack rules MUST be bound to a requirement.\n")
	testutil.WriteFile(t, root, "quill.toml", profiles.Format(t, config))

	engine, err := New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := engine.Check(context.Background(), CheckOptions{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	var echoed string
	for _, entry := range result.Rules {
		if entry.Rule.ID != "extpack/inspect" {
			continue
		}
		for _, diag := range entry.Execution.Diagnostics {
			echoed = diag.Message
		}
	}
	if echoed == "" {
		t.Fatal("expected the inspect rule to echo the request configuration")
	}
	if !strings.Contains(echoed, "database/sql") {
		t.Fatalf("expected external config to reach the runtime, got %q", echoed)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// stageExternalPack compiles the shared packhelper binary into a local Pack source directory
// beneath the repository root so the engine loads it as a real external Pack.
func stageExternalPack(t *testing.T, root string) {
	t.Helper()

	helperSource := filepath.Join(
		testutil.RepositoryRoot(t),
		"internal", "execution", "drivers", "testdata", "packhelper", "main.go",
	)
	binary := filepath.Join(t.TempDir(), "packhelper")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	command := exec.Command("go", "build", "-o", binary, helperSource)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("compile pack helper: %v\n%s", err, output)
	}

	contents, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	testutil.WriteExecutable(t, root, ".quill/packs/extpack/bin/packhelper", string(contents))
	testutil.WriteFile(t, root, ".quill/packs/extpack/pack.toml", externalAcceptanceManifest)
}

const externalAcceptanceManifest = `
schema_version = 1

[pack]
id = "extpack"
name = "External Test Pack"
version = "0.1.0"
quill_protocol = "quill-pack-v1"

[runtime]
command = "bin/packhelper"
timeout = "10s"

[[rules]]
id = "extpack/forbidden-import"
name = "Forbidden import"
check = "diagnostic"

[[rules]]
id = "extpack/marker"
name = "Execution marker"
check = "marker"

[[rules]]
id = "extpack/inspect"
name = "Inspect config"
check = "inspect-request"
`
