package drivers

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

func TestBuiltinPackExecutionKindsHaveDrivers(t *testing.T) {
	registry, err := builtin.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	checkers := CheckDrivers()
	fixers := FixDrivers()
	for _, rule := range registry.Rules() {
		if _, found := checkers[rule.Check.Kind]; !found {
			t.Fatalf("rule %q uses driver %q without a checker driver", rule.ID, rule.Check.Kind)
		}

		if rule.Fix.Empty() {
			continue
		}

		if _, found := fixers[rule.Fix.Kind]; !found {
			t.Fatalf("rule %q uses driver %q without a fixer driver", rule.ID, rule.Fix.Kind)
		}
	}
}

func TestMergeDriversRejectsDuplicateExecutionKind(t *testing.T) {
	driver := CheckDrivers()[contract.ExecutionToolchain]
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate execution kind to panic")
		}
	}()

	mergeDrivers(
		runner.DriverRegistry{contract.ExecutionToolchain: driver},
		runner.DriverRegistry{contract.ExecutionToolchain: driver},
	)
}
