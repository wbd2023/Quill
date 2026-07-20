package engine

import (
	"testing"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/pack/shipped/tool"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
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

	context := execution.RunContext{
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
