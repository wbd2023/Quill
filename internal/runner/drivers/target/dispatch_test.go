package target

import (
	"testing"

	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/golang/check"
)

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
