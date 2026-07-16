package pack

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Registry ------------------------------------------ */

func TestRegistryRejectsDuplicateRuleIDs(t *testing.T) {
	err := validateRegistry(buildRegistry([]Definition{
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
	}))
	if err == nil {
		t.Fatal("expected duplicate rule id to be rejected")
	}
}

func TestRegistryRejectsMissingCheckExecution(t *testing.T) {
	err := validateRegistry(buildRegistry([]Definition{
		{
			ID:   "broken",
			Name: "broken",
			Rules: []style.RuleDefinition{
				{ID: "missing/driver", Name: "missing driver"},
			},
		},
	}))
	if err == nil {
		t.Fatal("expected missing driver to be rejected")
	}
}

func TestRegistryRejectsBlankToolName(t *testing.T) {
	err := validateRegistry(buildRegistry([]Definition{
		{
			ID:   "one",
			Name: "one",
			Tools: []toolchain.Capability{
				{ID: "tool", Name: "", Command: "tool"},
			},
		},
	}))
	if err == nil {
		t.Fatal("expected blank tool name to be rejected")
	}
}

func TestRegistryRejectsDuplicatePackFileSets(t *testing.T) {
	err := validateRegistry(buildRegistry([]Definition{
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
	}))
	if err == nil {
		t.Fatal("expected duplicate pack file set to be rejected")
	}
}

func TestCatalogRegistryLoadsRegisteredPack(t *testing.T) {
	registry, err := NewCatalog(testPack("custom")).Registry([]string{"custom"})
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
	_, err := NewCatalog(testPack("duplicate"), testPack("duplicate")).Registry(nil)
	if err == nil {
		t.Fatal("expected duplicate pack id to fail")
	}
}

func TestCatalogPacksReturnIndependentCopies(t *testing.T) {
	catalog := NewCatalog(Definition{
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

func TestRegistryRulesReturnIndependentDefinitions(t *testing.T) {
	registry := buildRegistry([]Definition{
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

	rules := registry.Rules()
	execution := rules[0].Check.(style.FileCommandExecution)
	execution.Arguments[0] = "-changed"

	rules = registry.Rules()
	execution = rules[0].Check.(style.FileCommandExecution)
	if got := execution.Arguments[0]; got != "-w" {
		t.Fatalf("registry rule argument = %q, want -w", got)
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
