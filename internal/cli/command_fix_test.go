package cli

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
)

func TestFixableRulesUsePackFixes(t *testing.T) {
	rules := []contract.Rule{
		{
			ID: "go/lint",
			Fix: contract.ExecutionSpec{
				Kind: contract.ExecutorTargetCommand,
				Detail: contract.TargetCommandExecution{
					ToolIDs:  []string{builtin.ToolGo},
					Action:   builtin.TargetActionGoFormat,
					Language: builtin.LanguageGo,
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
				ScopeRoots: map[contract.Scope][]string{
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
