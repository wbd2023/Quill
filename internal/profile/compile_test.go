package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------------- Compile ------------------------------------------ */

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

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.FileSets = nil
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown file set") {
		t.Fatalf("expected unknown file set compile error, got %v", err)
	}
}

func TestCompileRejectsUnknownLanguageBackendBinding(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Language.Backends = nil
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown language backend") {
		t.Fatalf("expected unknown language backend compile error, got %v", err)
	}
}

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

func TestCompileRejectsExecutorSpecFieldMismatch(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Rules = []policy.RuleBinding{
		{
			RuleID:         "test/bad-file-command",
			Level:          contract.LevelRequired,
			Scope:          contract.Scope("all"),
			RequirementIDs: []string{"0.4.pin-go"},
		},
	}
	config.Tools = []policy.ToolPin{{ID: rulepack.ToolShfmt, Version: "v3.12.0"}}
	definitions := contract.Definitions{
		Tools: []contract.Tool{{ID: rulepack.ToolShfmt}},
		Rules: []contract.RuleDefinition{
			{
				ID: "test/bad-file-command",
				Spec: contract.ExecutionSpec{
					Kind: rulepack.ExecutorFileCommand,
					Detail: contract.FileCommandExecution{
						ToolID: rulepack.ToolShfmt,
					},
				},
			},
		},
	}

	_, err = Compile(config, definitions)
	if err == nil || !strings.Contains(err.Error(), "file-command spec must define a file set") {
		t.Fatalf("expected file-command shape error, got %v", err)
	}
}

func TestCompileRejectsUnknownRuleBindingPathClass(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	delete(config.Paths.Classes, "go_source")
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown path class") {
		t.Fatalf("expected missing path class compile error, got %v", err)
	}
}

func TestCompileRequiresActiveToolPins(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Tools = config.Tools[1:]
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "missing a tool pin") {
		t.Fatalf("expected missing tool pin error, got %v", err)
	}
}

func TestCompileRejectsUnknownToolPins(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Tools = append(config.Tools, policy.ToolPin{ID: "unknown", Version: "1.0.0"})
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected unknown tool pin error, got %v", err)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func effectiveRuleByID(
	effective contract.EffectiveConfig,
	ruleID string,
) (rule contract.Rule, found bool) {
	for _, rule := range effective.Rules {
		if rule.ID == ruleID {
			return rule, true
		}
	}

	return contract.Rule{}, false
}

func replaceRuleBindingBackends(
	config *policy.Config,
	ruleID string,
	backends []string,
) {
	for index := range config.Rules {
		if config.Rules[index].RuleID != ruleID {
			continue
		}

		config.Rules[index].Backends = append([]string{}, backends...)
		return
	}
}
