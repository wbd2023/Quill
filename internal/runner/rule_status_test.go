package runner

import (
	"errors"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
)

func TestCheckStatusRequiredViolationsFail(t *testing.T) {
	rule := contract.Rule{Level: contract.LevelRequired}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != report.CheckStatusFail {
		t.Fatalf("expected required violation to fail, got %q", status)
	}
}

func TestCheckStatusRecommendationsWarnByDefault(t *testing.T) {
	rule := contract.Rule{Level: contract.LevelRecommendation}

	status := CheckStatus(rule, errors.New("violation"), false)
	if status != report.CheckStatusWarn {
		t.Fatalf("expected recommendation violation to warn, got %q", status)
	}
}

func TestCheckStatusStrictRecommendationsFail(t *testing.T) {
	rule := contract.Rule{Level: contract.LevelRecommendation}

	status := CheckStatus(rule, errors.New("violation"), true)
	if status != report.CheckStatusFail {
		t.Fatalf("expected strict recommendation violation to fail, got %q", status)
	}
}

func TestCheckStatusBlockedRulesSkip(t *testing.T) {
	rule := contract.Rule{Level: contract.LevelRequired}

	status := CheckStatus(rule, errRuleBlocked, false)
	if status != report.CheckStatusSkip {
		t.Fatalf("expected blocked rule to skip, got %q", status)
	}
}
