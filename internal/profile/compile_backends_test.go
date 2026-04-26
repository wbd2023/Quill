package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------- Rule Binding Checks ------------------------------------ */

func TestCompileBindsLanguageBackendsFromRuleBinding(t *testing.T) {
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

	rule, found := effectiveRuleByID(effective, "go/lint")
	if !found {
		t.Fatal("go/lint rule missing")
	}

	if strings.Join(rule.Spec.Backends(), ",") != "application_go,tooling_go" {
		t.Fatalf("go/lint backends = %v", rule.Spec.Backends())
	}

	if strings.Join(rule.FixSpec.Backends(), ",") != "application_go,tooling_go" {
		t.Fatalf("go/lint fix backends = %v", rule.FixSpec.Backends())
	}
}

func TestCompileBindsPathClassesFromRuleBinding(t *testing.T) {
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

	rule, found := effectiveRuleByID(effective, "go/errors")
	if !found {
		t.Fatal("go/errors rule missing")
	}

	if strings.Join(rule.PathClasses, ",") != "go_source,concrete_infra,domain_errors" {
		t.Fatalf("go/errors path classes = %v", rule.PathClasses)
	}
}

func TestCompileRejectsDuplicateBackendBinding(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	replaceRuleBindingBackends(&config, "go/lint", []string{"application_go", "application_go"})
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "duplicates backend") {
		t.Fatalf("expected duplicate backend compile error, got %v", err)
	}
}

func TestCompileRejectsBackendsOnNonLanguageRule(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	replaceRuleBindingBackends(&config, "text/ascii", []string{"application_go"})
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unexpected backends") {
		t.Fatalf("expected unexpected backend compile error, got %v", err)
	}
}
