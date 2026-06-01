package runner

import (
	"errors"
	"testing"

	"ciphera/tools/internal/contract"
)

func TestCheckStatusRequiredViolationsFail(t *testing.T) {
	rule := contract.Rule{Enforcement: contract.EnforcementRequired}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != contract.CheckStatusFail {
		t.Fatalf("expected required violation to fail, got %q", status)
	}
}

func TestCheckStatusRecommendationsWarnByDefault(t *testing.T) {
	rule := contract.Rule{Enforcement: contract.EnforcementRecommendation}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != contract.CheckStatusWarn {
		t.Fatalf("expected recommendation violation to warn, got %q", status)
	}
}

func TestCheckStatusStrictRecommendationsFail(t *testing.T) {
	rule := contract.Rule{Enforcement: contract.EnforcementRecommendation}

	status := CheckStatus(rule, errors.New("violation"), true)
	if status != contract.CheckStatusFail {
		t.Fatalf("expected strict recommendation violation to fail, got %q", status)
	}
}

func TestCheckStatusBlockedRulesSkip(t *testing.T) {
	rule := contract.Rule{Enforcement: contract.EnforcementRequired}

	status := CheckStatus(rule, errRuleBlocked, false)
	if status != contract.CheckStatusSkip {
		t.Fatalf("expected blocked rule to skip, got %q", status)
	}
}
