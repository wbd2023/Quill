package drivers

import (
	"testing"

	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/golang/check"
)

func TestBuiltinPackExecutorsHaveDrivers(t *testing.T) {
	registry, err := builtin.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checkers := Checkers()
	fixers := Fixers()
	for _, rule := range registry.Rules() {
		if _, found := checkers[rule.Check.Kind]; !found {
			t.Fatalf("rule %q uses executor %q without a checker driver", rule.ID, rule.Check.Kind)
		}

		if rule.Fix.Empty() {
			continue
		}

		if _, found := fixers[rule.Fix.Kind]; !found {
			t.Fatalf("rule %q uses executor %q without a fixer driver", rule.ID, rule.Fix.Kind)
		}
	}
}

func TestBuiltinPackScannersHaveDrivers(t *testing.T) {
	registry, err := builtin.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	scanners := repositoryScanners()
	for _, rule := range registry.Rules() {
		execution, found := rule.Check.RepositoryScanExecution()
		if !found {
			continue
		}

		if _, found := scanners[execution.Scanner]; !found {
			t.Fatalf(
				"rule %q uses scanner %q without a driver",
				rule.ID,
				execution.Scanner,
			)
		}
	}
}

func TestBuiltinPackGoChecksHaveDispatch(t *testing.T) {
	registry, err := builtin.DefaultRegistry([]string{builtin.PackGo})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checks := map[string]bool{}
	for _, checkID := range check.IDs() {
		checks[checkID] = true
	}

	for _, rule := range registry.Rules() {
		execution, found := rule.Check.TargetCheckExecution()
		if !found {
			continue
		}

		if !checks[execution.Check] {
			t.Fatalf(
				"rule %q uses Go check %q without dispatch",
				rule.ID,
				execution.Check,
			)
		}
	}
}
