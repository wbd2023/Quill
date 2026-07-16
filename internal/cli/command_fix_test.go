package cli

import (
	"testing"

	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func TestFixableRulesUsePackFixes(t *testing.T) {
	rules := []style.Rule{
		{
			ID: "go/lint",
			Fix: style.TargetCommandJob{
				ToolIDs:  []string{tool.Go},
				Action:   golang.TargetActionGoFormat,
				Language: golang.Language,
			},
			Scope: style.Scope("tools"),
		},
		{
			ID:    "security/secrets",
			Scope: style.Scope("all"),
		},
	}

	context := runner.Context{
		Scope: style.Scope("tools"),
		Profile: policy.Config{
			Repository: policy.RepositoryConfig{
				ScopeRoots: map[style.Scope][]string{
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
