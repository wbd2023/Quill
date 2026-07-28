package pack

import (
	"testing"

	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------------ Registry ------------------------------------------ */

func TestRegistryRejectsDuplicateRuleIDs(t *testing.T) {
	registry, err := buildRegistry(nil, []Definition{
		{
			ID:   "one",
			Name: "one",
			Rules: []style.RuleDefinition{
				{
					ID:   "duplicate",
					Name: "first",
					Check: style.RepositoryScanExecution{
						Scanner: "test",
					},
				},
			},
		},
		{
			ID:   "two",
			Name: "two",
			Rules: []style.RuleDefinition{
				{
					ID:   "duplicate",
					Name: "second",
					Check: style.RepositoryScanExecution{
						Scanner: "test",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildRegistry: %v", err)
	}

	if err = validateRegistry(registry); err == nil {
		t.Fatal("expected duplicate rule id to be rejected")
	}
}

func TestRegistryRejectsMissingCheckExecution(t *testing.T) {
	registry, err := buildRegistry(nil, []Definition{
		{
			ID:   "broken",
			Name: "broken",
			Rules: []style.RuleDefinition{
				{ID: "missing/driver", Name: "missing driver"},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildRegistry: %v", err)
	}

	if err = validateRegistry(registry); err == nil {
		t.Fatal("expected missing driver to be rejected")
	}
}

func TestRegistryRejectsDuplicatePackFileSets(t *testing.T) {
	registry, err := buildRegistry(nil, []Definition{
		{
			ID:   "one",
			Name: "one",
			FileSets: policy.FileSets{
				{Name: "source"},
			},
		},
		{
			ID:   "two",
			Name: "two",
			FileSets: policy.FileSets{
				{Name: "source"},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildRegistry: %v", err)
	}

	if err = validateRegistry(registry); err == nil {
		t.Fatal("expected duplicate pack file set to be rejected")
	}
}

func TestRegistryRulesReturnIndependentDefinitions(t *testing.T) {
	registry, err := buildRegistry(nil, []Definition{
		{
			ID:   "custom",
			Name: "Custom",
			Rules: []style.RuleDefinition{
				{
					ID:   "custom/rule",
					Name: "Custom rule",
					Check: style.FileCommandExecution{
						Arguments: []string{"-w"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildRegistry: %v", err)
	}

	rules := registry.Rules()
	execution := rules[0].Check.(style.FileCommandExecution)
	execution.Arguments[0] = "-changed"

	rules = registry.Rules()
	execution = rules[0].Check.(style.FileCommandExecution)
	if got := execution.Arguments[0]; got != "-w" {
		t.Fatalf("registry rule argument = %q, want -w", got)
	}
}

/* ---------------------------------- Catalogue Tool Ownership ---------------------------------- */

func TestCatalogRejectsDuplicateToolIDs(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "go", Name: "Go", Command: "go"},
		{ID: "go", Name: "Go Again", Command: "go"},
	}

	_, err := NewCatalog(tools, testPack("custom")).Registry(nil)
	if err == nil {
		t.Fatal("expected duplicate tool id to be rejected before loss")
	}
}

func TestCatalogRejectsBlankToolName(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "tool", Name: "", Command: "tool"},
	}

	_, err := NewCatalog(tools, testPack("custom")).Registry(nil)
	if err == nil {
		t.Fatal("expected blank tool name to be rejected")
	}
}

func TestRegistryRejectsUnknownToolReference(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "go", Name: "Go", Command: "go"},
	}

	packWithUnknownTool := Definition{
		ID:      "custom",
		Name:    "Custom",
		ToolIDs: []string{"go", "ghost"},
		Rules: []style.RuleDefinition{
			{
				ID:   "custom/rule",
				Name: "Custom rule",
				Check: style.ToolchainExecution{
					ToolIDs: []string{"go", "ghost"},
				},
			},
		},
	}

	if _, err := NewCatalog(tools, packWithUnknownTool).Registry(nil); err == nil {
		t.Fatal("expected unknown tool reference to be rejected")
	}
}

func TestRegistryResolvesSharedToolOnce(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "go", Name: "Go", Command: "go"},
	}

	twoPacks := []Definition{
		{
			ID:      "alpha",
			Name:    "Alpha",
			ToolIDs: []string{"go"},
			Rules: []style.RuleDefinition{
				{
					ID:    "alpha/rule",
					Name:  "Alpha rule",
					Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
				},
			},
		},
		{
			ID:      "beta",
			Name:    "Beta",
			ToolIDs: []string{"go"},
			Rules: []style.RuleDefinition{
				{
					ID:    "beta/rule",
					Name:  "Beta rule",
					Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
				},
			},
		},
	}

	registry, err := NewCatalog(tools, twoPacks...).Registry(nil)
	if err != nil {
		t.Fatalf("Registry: %v", err)
	}

	if len(registry.ToolCapabilities()) != 1 {
		t.Fatalf(
			"expected shared tool resolved once, got %d capabilities",
			len(registry.ToolCapabilities()),
		)
	}
}

// A Pack's rule may only reference Tools the Pack itself declares. Even when the aggregate
// catalogue contains a Tool (because another Pack declares it), a Pack borrowing that Tool for its
// own rule is rejected at assembly.
func TestRegistryRejectsCrossPackToolReference(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "go", Name: "Go", Command: "go"},
	}

	leakyPack := Definition{
		ID:      "alpha",
		Name:    "Alpha",
		ToolIDs: nil,
		Rules: []style.RuleDefinition{
			{
				ID:    "alpha/versions",
				Name:  "Alpha versions",
				Group: "project",
				Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
			},
		},
	}

	owningPack := Definition{
		ID:      "beta",
		Name:    "Beta",
		ToolIDs: []string{"go"},
		Rules: []style.RuleDefinition{
			{
				ID:    "beta/versions",
				Name:  "Beta versions",
				Group: "project",
				Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
			},
		},
	}

	if _, err := NewCatalog(tools, leakyPack, owningPack).Registry(nil); err == nil {
		t.Fatal("expected cross-pack tool reference to be rejected")
	}
}

func TestRegistryAcceptsInPackToolReference(t *testing.T) {
	tools := []toolchain.Capability{
		{ID: "go", Name: "Go", Command: "go"},
	}

	owningPack := Definition{
		ID:      "beta",
		Name:    "Beta",
		ToolIDs: []string{"go"},
		Rules: []style.RuleDefinition{
			{
				ID:    "beta/versions",
				Name:  "Beta versions",
				Group: "project",
				Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
			},
		},
	}

	if _, err := NewCatalog(tools, owningPack).Registry(nil); err != nil {
		t.Fatalf("expected in-pack tool reference to validate, got: %v", err)
	}
}

/* ----------------------------------------- Provenance ----------------------------------------- */

func TestRegistryStampsPackIDOntoRules(t *testing.T) {
	registry, err := buildRegistry(nil, []Definition{testPack("provenance")})
	if err != nil {
		t.Fatalf("buildRegistry: %v", err)
	}

	rules := registry.Rules()
	if len(rules) != 1 {
		t.Fatalf("rules = %d, want 1", len(rules))
	}

	if got := rules[0].PackID; got != "provenance" {
		t.Fatalf("rule PackID = %q, want provenance", got)
	}

	execution, ok := rules[0].Check.(style.RepositoryScanExecution)
	if !ok {
		t.Fatalf("expected RepositoryScanExecution, got %T", rules[0].Check)
	}

	if got := execution.PackID; got != "provenance" {
		t.Fatalf("execution PackID = %q, want provenance", got)
	}
}

func TestRegistryToolCapabilitiesAreDefensiveCopies(t *testing.T) {
	tool := toolchain.Capability{
		ID:      "shellcheck",
		Name:    "shellcheck",
		Command: "shellcheck",
		Install: toolchain.GitHubInstall{Platforms: map[string]string{"linux/amd64": "x"}},
	}

	packReferencingTool := Definition{
		ID:      "custom",
		Name:    "Custom",
		ToolIDs: []string{"shellcheck"},
		Rules: []style.RuleDefinition{
			{
				ID:    "custom/rule",
				Name:  "Custom rule",
				Group: "external_tools",
				Check: style.ToolchainExecution{ToolIDs: []string{"shellcheck"}},
			},
		},
	}

	registry, err := NewCatalog([]toolchain.Capability{tool}, packReferencingTool).Registry(nil)
	if err != nil {
		t.Fatalf("Registry: %v", err)
	}

	capabilities := registry.ToolCapabilities()
	capabilities[0].Install.(toolchain.GitHubInstall).Platforms["linux/amd64"] = "mutated"

	again := registry.ToolCapabilities()
	if got := again[0].Install.(toolchain.GitHubInstall).Platforms["linux/amd64"]; got != "x" {
		t.Fatalf("registry capability mutated via returned copy: %q", got)
	}
}

/* ------------------------------------------- Catalog ------------------------------------------ */

func TestCatalogRegistryLoadsRegisteredPack(t *testing.T) {
	registry, err := NewCatalog(nil, testPack("custom")).Registry([]string{"custom"})
	if err != nil {
		t.Fatalf("Registry: %v", err)
	}

	if len(registry.Packs()) != 1 {
		t.Fatalf("packs = %d, want 1", len(registry.Packs()))
	}

	if len(registry.Rules()) != 1 {
		t.Fatalf("rules = %d, want 1", len(registry.Rules()))
	}
}

func TestCatalogRegistryRejectsDuplicatePackIDs(t *testing.T) {
	_, err := NewCatalog(nil, testPack("duplicate"), testPack("duplicate")).Registry(nil)
	if err == nil {
		t.Fatal("expected duplicate pack id to fail")
	}
}

func TestCatalogPacksReturnIndependentCopies(t *testing.T) {
	catalog := NewCatalog(nil, Definition{
		ID:   "custom",
		Name: "Custom",
		FileSets: policy.FileSets{
			{
				Name: "source",
				Include: policy.FileSetInclude{
					Extensions: []string{".go"},
				},
			},
		},
	})

	packs := catalog.Packs()
	packs[0].FileSets[0].Include.Extensions[0] = ".txt"

	packs = catalog.Packs()
	if got := packs[0].FileSets[0].Include.Extensions[0]; got != ".go" {
		t.Fatalf("catalog pack file set extension = %q, want .go", got)
	}
}

func TestCatalogToolsReturnIndependentCopies(t *testing.T) {
	tools := []toolchain.Capability{
		{
			ID:      "shellcheck",
			Name:    "shellcheck",
			Command: "shellcheck",
			Install: toolchain.GitHubInstall{Platforms: map[string]string{"linux/amd64": "x"}},
		},
	}

	catalog := NewCatalog(tools, testPack("custom"))
	returned := catalog.Tools()
	returned[0].Install.(toolchain.GitHubInstall).Platforms["linux/amd64"] = "mutated"

	again := catalog.Tools()
	if got := again[0].Install.(toolchain.GitHubInstall).Platforms["linux/amd64"]; got != "x" {
		t.Fatalf("catalog tool mutated via returned copy: %q", got)
	}
}

func testPack(id string) (definition Definition) {
	return Definition{
		ID:   id,
		Name: id,
		Rules: []style.RuleDefinition{
			{
				ID:   id + "/rule",
				Name: id + " rule",
				Check: style.RepositoryScanExecution{
					Scanner: "test",
				},
			},
		},
	}
}
