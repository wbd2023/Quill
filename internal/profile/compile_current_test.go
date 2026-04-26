package profile

import (
	"testing"

	"ciphera/tools/internal/rulepack"
)

func TestCurrentProfileCompilesEnabledRulePacks(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(effective.Rules) != len(registry.Rules()) {
		t.Fatalf("expected %d effective rules, got %d", len(registry.Rules()), len(effective.Rules))
	}
}
