package runner

import (
	"errors"
	"testing"

	"ciphera/tools/internal/style"
)

func TestCheckStatusRequiredViolationsFail(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != style.CheckStatusFail {
		t.Fatalf("expected required violation to fail, got %q", status)
	}
}

func TestCheckStatusRecommendationsWarnByDefault(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRecommendation}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != style.CheckStatusWarn {
		t.Fatalf("expected recommendation violation to warn, got %q", status)
	}
}

func TestCheckStatusStrictRecommendationsFail(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRecommendation}

	status := CheckStatus(rule, errors.New("violation"), true)
	if status != style.CheckStatusFail {
		t.Fatalf("expected strict recommendation violation to fail, got %q", status)
	}
}

func TestCheckStatusBlockedRulesSkip(t *testing.T) {
	rule := style.Rule{Enforcement: style.EnforcementRequired}

	status := CheckStatus(rule, errRuleBlocked, false)
	if status != style.CheckStatusSkip {
		t.Fatalf("expected blocked rule to skip, got %q", status)
	}
}
