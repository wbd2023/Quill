package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------------- Compile ------------------------------------------ */

func TestCurrentProfileCompilesEnabledRulePacks(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(effective.Rules) != len(registry.Rules()) {
		t.Fatalf("expected %d effective rules, got %d", len(registry.Rules()), len(effective.Rules))
	}
}

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	policy.FileSets = nil
	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = policy.Compile(registry)
	if err == nil || !strings.Contains(err.Error(), "unknown file set") {
		t.Fatalf("expected unknown file set compile error, got %v", err)
	}
}

func TestCompileRejectsUnknownLanguageBackendBinding(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	policy.Language.Backends = nil
	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = policy.Compile(registry)
	if err == nil || !strings.Contains(err.Error(), "unknown language backend") {
		t.Fatalf("expected unknown language backend compile error, got %v", err)
	}
}

func TestCompileRejectsMissingGoParameterVocabulary(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	policy.Naming.GoParameters.SecretNames = nil
	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = policy.Compile(registry)
	if err == nil || !strings.Contains(err.Error(), "naming.go_parameters.secret_names") {
		t.Fatalf("expected missing Go parameter vocabulary error, got %v", err)
	}
}

func TestCompileRejectsMissingRulePackPathClass(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	delete(policy.Paths.Classes, rulepack.PathClassApp)
	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = policy.Compile(registry)
	if err == nil || !strings.Contains(err.Error(), "requires paths.app") {
		t.Fatalf("expected missing path class compile error, got %v", err)
	}
}
