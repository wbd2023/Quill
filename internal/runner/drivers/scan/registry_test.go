package scan

import (
	"testing"

	"ciphera/tools/internal/pack/builtin"
)

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
