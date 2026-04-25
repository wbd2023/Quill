package cli

import (
	"testing"

	"ciphera/tools/internal/contract"
)

func TestFixableRulesUsesRulePackFixSpecs(t *testing.T) {
	rules := []contract.Rule{
		{
			RuleDefinition: contract.RuleDefinition{
				ID: "go/lint-tools",
				FixSpec: contract.ExecutionSpec{
					Executor: contract.ExecutorGoFormat,
				},
			},
			Scope: contract.ScopeTools,
		},
		{
			RuleDefinition: contract.RuleDefinition{
				ID: "repo/secrets",
			},
			Scope: contract.ScopeAll,
		},
	}

	selected := fixableRules(rules, contract.ScopeTools)
	if len(selected) != 1 || selected[0].ID != "go/lint-tools" {
		t.Fatalf("fixableRules = %v", selected)
	}
}
