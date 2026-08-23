package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

func TestExplainReturnsActiveRuleWithoutToolInspection(t *testing.T) {
	t.Parallel()

	instance, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := instance.Explain(context.Background(), "go/lint")
	if err != nil {
		t.Fatalf("Explain: %v", err)
	}
	if result.Rule.ID != "go/lint" || !result.Rule.Active {
		t.Fatalf("Explain result = %#v", result.Rule)
	}
}

func TestExplainRejectsUnknownAndInactiveRulesAsArguments(t *testing.T) {
	t.Parallel()

	instance, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for _, ruleID := range []string{"does/not-exist", "go/file-shape"} {
		t.Run(ruleID, func(t *testing.T) {
			_, err := instance.Explain(context.Background(), ruleID)
			var argumentError *ArgumentError
			if !errors.As(err, &argumentError) {
				t.Fatalf("Explain(%q) error = %v, want ArgumentError", ruleID, err)
			}
		})
	}
}
