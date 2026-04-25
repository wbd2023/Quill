package rulepack

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
)

/* ------------------------------------------- Tooling ------------------------------------------ */

func TestRegistryToolsHaveUniqueIDs(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	seenIDs := make(map[string]bool)
	for _, tool := range registry.Tools() {
		if seenIDs[tool.ID] {
			t.Fatalf("duplicate tool ID: %s", tool.ID)
		}

		seenIDs[tool.ID] = true
	}
}

func TestRegistryToolsUseSupportedInstallStrategies(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, tool := range registry.Tools() {
		switch tool.InstallKind {
		case contract.ToolInstallNone,
			contract.ToolInstallGoBinary,
			contract.ToolInstallNodePackage,
			contract.ToolInstallShellcheckArchive:
		default:
			t.Fatalf("tool %q uses unsupported install strategy %q", tool.ID, tool.InstallKind)
		}

		switch tool.InstallKind {
		case contract.ToolInstallGoBinary, contract.ToolInstallNodePackage:
			if tool.InstallSource == "" {
				t.Fatalf("tool %q must define an install source", tool.ID)
			}
		case contract.ToolInstallNone, contract.ToolInstallShellcheckArchive:
			if tool.InstallSource != "" {
				t.Fatalf("tool %q must not define an install source", tool.ID)
			}
		}
	}
}

/* ---------------------------------------- Version Pins ---------------------------------------- */

func TestPinnedGoVersionMatchesModuleFiles(t *testing.T) {
	goDirectivePattern := regexp.MustCompile(`(?m)^go ([0-9]+\.[0-9]+(?:\.[0-9]+)?)$`)
	goTool := toolByID(t, contract.ToolGo)

	rootModule := readRepoFile(t, "go.mod")
	styleModule := readRepoFile(t, filepath.Join("tools", "go.mod"))

	for _, contents := range []string{rootModule, styleModule} {
		matches := goDirectivePattern.FindStringSubmatch(contents)
		if len(matches) != 2 {
			t.Fatalf("could not find go directive in module contents:\n%s", contents)
		}

		if matches[1] != goTool.PinnedVersion {
			t.Fatalf(
				"go directive %q does not match pinned version %q",
				matches[1],
				goTool.PinnedVersion,
			)
		}
	}
}

func TestPinnedGoimportsVersionMatchesStyleModule(t *testing.T) {
	requireLinePattern := regexp.MustCompile(
		`(?m)^\s*golang\.org/x/tools (v[0-9]+\.[0-9]+\.[0-9]+)(?:$| // indirect$)`,
	)
	goimportsTool := toolByID(t, contract.ToolGoimports)

	styleModule := readRepoFile(t, filepath.Join("tools", "go.mod"))
	matches := requireLinePattern.FindStringSubmatch(styleModule)
	if len(matches) != 2 {
		t.Fatalf("could not find golang.org/x/tools requirement in tools/go.mod")
	}

	if matches[1] != goimportsTool.PinnedVersion {
		t.Fatalf(
			"tools/go.mod pins golang.org/x/tools at %q, want %q",
			matches[1],
			goimportsTool.PinnedVersion,
		)
	}
}

/* -------------------------------------- Repository Files -------------------------------------- */

func toolByID(t *testing.T, toolID string) (tool contract.Tool) {
	t.Helper()

	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	tool, found := registry.ToolByID(toolID)
	if !found {
		t.Fatalf("missing %s tool in registry", toolID)
	}

	return tool
}

func readRepoFile(t *testing.T, relativePath string) (contents string) {
	t.Helper()

	path := filepath.Join(fixtures.RepoRoot(t), filepath.FromSlash(relativePath))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", relativePath, err)
	}

	return strings.TrimSpace(string(data))
}
