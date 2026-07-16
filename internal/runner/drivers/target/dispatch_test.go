package target

import (
	"testing"

	"ciphera/tools/internal/checks/golang/check"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/style"
)

func TestShippedPackGoChecksHaveDispatch(t *testing.T) {
	registry, err := shipped.DefaultRegistry([]string{golang.PackID})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checks := map[string]bool{}
	for _, checkID := range check.IDs() {
		checks[checkID] = true
	}

	for _, rule := range registry.Rules() {
		execution, found := rule.Check.(style.TargetCheckTemplate)
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
