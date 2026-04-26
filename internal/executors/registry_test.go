package executors

import (
	"testing"

	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/golang"
)

func TestRulepackExecutorsHaveBindings(t *testing.T) {
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checkers := Checkers()
	fixers := Fixers()
	for _, rule := range registry.Rules() {
		if _, found := checkers[rule.Spec.Kind]; !found {
			t.Fatalf("rule %q uses executor %q without a checker binding", rule.ID, rule.Spec.Kind)
		}

		if rule.FixSpec.Empty() {
			continue
		}

		if _, found := fixers[rule.FixSpec.Kind]; !found {
			t.Fatalf("rule %q uses executor %q without a fixer binding", rule.ID, rule.FixSpec.Kind)
		}
	}
}

func TestRulepackScannersHaveBindings(t *testing.T) {
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	scanners := repositoryScanners()
	for _, rule := range registry.Rules() {
		detail, found := rule.Spec.RepositoryScanExecution()
		if !found {
			continue
		}

		if _, found := scanners[detail.Scanner]; !found {
			t.Fatalf(
				"rule %q uses scanner %q without a binding",
				rule.ID,
				detail.Scanner,
			)
		}
	}
}

func TestRulepackGoChecksHaveDispatch(t *testing.T) {
	registry, err := rulepack.DefaultRegistry([]string{rulepack.PackGo})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checks := map[string]bool{}
	for _, checkID := range golang.CheckIDs() {
		checks[checkID] = true
	}

	for _, rule := range registry.Rules() {
		detail, found := rule.Spec.BackendCheckExecution()
		if !found {
			continue
		}

		if !checks[detail.Check] {
			t.Fatalf(
				"rule %q uses Go check %q without dispatch",
				rule.ID,
				detail.Check,
			)
		}
	}
}
