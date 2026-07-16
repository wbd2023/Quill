package engine

import (
	"testing"

	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func TestSelectRulesForFixFiltersByScopeAndFixPresence(t *testing.T) {
	t.Parallel()

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

	selected := selectRulesForFix(rules, context)
	if len(selected) != 1 || selected[0].ID != "go/lint" {
		t.Fatalf("selectRulesForFix = %v", selected)
	}
}
