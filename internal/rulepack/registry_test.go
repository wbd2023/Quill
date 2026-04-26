package rulepack

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Registry ------------------------------------------ */

func TestDefaultRegistryLoadsEnabledRulePacks(t *testing.T) {
	registry, err := DefaultRegistry([]string{PackControl, PackText, PackMarkdown})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	if len(registry.Packs()) != 3 {
		t.Fatalf("expected 3 packs, got %d", len(registry.Packs()))
	}

	if len(registry.Rules()) == 0 {
		t.Fatal("expected enabled packs to register rules")
	}

	if _, found := registry.ToolByID(ToolMarkdownlint); !found {
		t.Fatal("expected markdown pack tool to be registered")
	}
}

func TestDefaultRegistryRejectsUnknownRulePack(t *testing.T) {
	if _, err := DefaultRegistry([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown rule pack to be rejected")
	}
}

func TestDefaultRegistryRejectsDuplicateRulePack(t *testing.T) {
	if _, err := DefaultRegistry([]string{PackGo, PackGo}); err == nil {
		t.Fatal("expected duplicate rule pack to be rejected")
	}
}

func TestRegistryRejectsDuplicateRuleIDs(t *testing.T) {
	err := validateRegistry(buildRegistry([]Pack{
		{
			ID:   "one",
			Name: "one",
			Rules: []RuleDefinition{
				{
					ID:   "duplicate",
					Name: "first",
					Spec: contract.ExecutionSpec{
						Kind: ExecutorRepositoryScan,
						Detail: contract.RepositoryScanExecution{
							Scanner: "test",
						},
					},
				},
			},
		},
		{
			ID:   "two",
			Name: "two",
			Rules: []RuleDefinition{
				{
					ID:   "duplicate",
					Name: "second",
					Spec: contract.ExecutionSpec{
						Kind: ExecutorRepositoryScan,
						Detail: contract.RepositoryScanExecution{
							Scanner: "test",
						},
					},
				},
			},
		},
	}))
	if err == nil {
		t.Fatal("expected duplicate rule id to be rejected")
	}
}

func TestRegistryRejectsMissingExecutor(t *testing.T) {
	err := validateRegistry(buildRegistry([]Pack{
		{
			ID:   "broken",
			Name: "broken",
			Rules: []RuleDefinition{
				{ID: "missing/executor", Name: "missing executor"},
			},
		},
	}))
	if err == nil {
		t.Fatal("expected missing executor to be rejected")
	}
}

func TestRegistryRejectsConflictingToolDefinitions(t *testing.T) {
	err := validateRegistry(buildRegistry([]Pack{
		{
			ID:   "one",
			Name: "one",
			Tools: []toolchain.Capability{
				{ID: "tool", Name: "tool", Command: "first"},
			},
		},
		{
			ID:   "two",
			Name: "two",
			Tools: []toolchain.Capability{
				{ID: "tool", Name: "tool", Command: "second"},
			},
		},
	}))
	if err == nil {
		t.Fatal("expected conflicting tool definitions to be rejected")
	}
}
