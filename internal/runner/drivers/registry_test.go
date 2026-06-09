package drivers

import (
	"testing"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func TestShippedPackExecutionKindsHaveDrivers(t *testing.T) {
	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	bindings := NewBindings()
	checkers := CheckDrivers(bindings)
	fixers := FixDrivers(bindings)
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
	driver := CheckDrivers(NewBindings())[style.ExecutionToolchain]
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate execution kind to panic")
		}
	}()

	mergeDrivers(
		runner.DriverRegistry{style.ExecutionToolchain: driver},
		runner.DriverRegistry{style.ExecutionToolchain: driver},
	)
}
