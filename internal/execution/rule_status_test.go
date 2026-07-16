package execution

import (
	"errors"
	"testing"

	"ciphera/tools/internal/style"
)

func TestCheckStatusRequiredViolationsFail(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}
	violations := style.ExecutionResult{
		Diagnostics: []style.Diagnostic{{Code: "test", Message: "violation"}},
	}

	status := CheckStatus(rule, violations, nil, false)
	if status != style.CheckStatusFail {
		t.Fatalf("expected required violation to fail, got %q", status)
	}
}

func TestCheckStatusRecommendationsWarnByDefault(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRecommendation}
	violations := style.ExecutionResult{
		Diagnostics: []style.Diagnostic{{Code: "test", Message: "violation"}},
	}

	status := CheckStatus(rule, violations, nil, false)
	if status != style.CheckStatusWarn {
		t.Fatalf("expected recommendation violation to warn, got %q", status)
	}
}

func TestCheckStatusStrictRecommendationsFail(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRecommendation}
	violations := style.ExecutionResult{
		Diagnostics: []style.Diagnostic{{Code: "test", Message: "violation"}},
	}

	status := CheckStatus(rule, violations, nil, true)
	if status != style.CheckStatusFail {
		t.Fatalf("expected strict recommendation violation to fail, got %q", status)
	}
}

func TestCheckStatusBlockedRulesSkip(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}

	status := CheckStatus(rule, style.ExecutionResult{}, errRuleBlocked, false)
	if status != style.CheckStatusSkip {
		t.Fatalf("expected blocked rule to skip, got %q", status)
	}
}

func TestCheckStatusOperationalErrorsError(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}

	status := CheckStatus(rule, style.ExecutionResult{}, errors.New("parse failed"), false)
	if status != style.CheckStatusError {
		t.Fatalf("expected operational error to error, got %q", status)
	}
}

func TestCheckStatusCleanResultsPass(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}

	status := CheckStatus(rule, style.ExecutionResult{}, nil, false)
	if status != style.CheckStatusPass {
		t.Fatalf("expected clean result to pass, got %q", status)
	}
}
