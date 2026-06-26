package drivers

import (
	"testing"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func TestShippedPackExecutionDetailsHaveDrivers(t *testing.T) {
	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	bindings := NewBindings()
	checkers := CheckDrivers(bindings)
	fixers := FixDrivers(bindings)
	for _, rule := range registry.Rules() {
		driver := driverForDetail(rule.Check.Detail, checkers)
		if driver == nil {
			t.Fatalf("rule %q check detail %T has no checker driver", rule.ID, rule.Check.Detail)
		}

		if rule.Fix.Empty() {
			continue
		}

		driver = driverForDetail(rule.Fix.Detail, fixers)
		if driver == nil {
			t.Fatalf("rule %q fix detail %T has no fixer driver", rule.ID, rule.Fix.Detail)
		}
	}
}

func driverForDetail(detail style.ExecutionDetail, set runner.DriverSet) (driver runner.Driver) {
	switch detail.(type) {

	case style.ToolchainExecution:
		return set.Toolchain

	case style.ProfileExecution:
		return set.Profile

	case style.FileCommandExecution:
		return set.FileCommand

	case style.TargetCommandExecution:
		return set.TargetCommand

	case style.TargetCheckExecution:
		return set.TargetCheck

	case style.RepositoryScanExecution:
		return set.RepositoryScan

	default:
		return nil
	}
}
