package cli

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func TestFixableRulesUsesRulePackFixSpecs(t *testing.T) {
	rules := []contract.Rule{
		{
			ID: "go/lint",
			FixSpec: contract.ExecutionSpec{
				Kind: rulepack.ExecutorBackendCommand,
				Detail: contract.BackendCommandExecution{
					ToolIDs:  []string{rulepack.ToolGo},
					Action:   rulepack.BackendActionGoFormat,
					Language: rulepack.LanguageGo,
				},
			},
			Scope: contract.Scope("tools"),
		},
		{
			ID:    "security/secrets",
			Scope: contract.Scope("all"),
		},
	}

	context := runner.Context{
		Scope: contract.Scope("tools"),
		Policy: policy.Config{
			Repository: policy.RepositoryConfig{
				Scopes: map[contract.Scope][]string{
					"all":   {"."},
					"tools": {"tools"},
				},
			},
		},
	}
	selected := fixableRules(rules, context)
	if len(selected) != 1 || selected[0].ID != "go/lint" {
		t.Fatalf("fixableRules = %v", selected)
	}
}
