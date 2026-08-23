package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
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
		Profile: profile.Profile{
			Repository: profile.RepositoryConfig{
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

func TestCheckRejectsCancelledContext(t *testing.T) {
	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	operationContext, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = engine.Check(operationContext, CheckOptions{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Check error = %v, want context.Canceled", err)
	}
}
