package drivers

import (
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/style"
)

func TestShippedPackExecutionDetailsHaveDrivers(t *testing.T) {
	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	bindings := NewBindings()
	checkers := Checks(bindings)
	fixers := Fixes(bindings)
	for _, rule := range registry.Rules() {
		driver := driverForDetail(rule.Check, checkers)
		if driver == nil {
			t.Fatalf("rule %q check detail %T has no checker driver", rule.ID, rule.Check)
		}

		if rule.Fix == nil {
			continue
		}

		driver = driverForDetail(rule.Fix, fixers)
		if driver == nil {
			t.Fatalf("rule %q fix detail %T has no fixer driver", rule.ID, rule.Fix)
		}
	}
}

func driverForDetail(detail style.Template, set execution.DriverSet) (driver execution.Driver) {
	switch detail.(type) {

	case style.ToolchainCheck:
		return set.Toolchain

	case style.ProfileCheck:
		return set.Profile

	case style.FileCommand:
		return set.FileCommand

	case style.TargetCommandTemplate:
		return set.TargetCommand

	case style.TargetCheckTemplate:
		return set.TargetCheck

	case style.RepositoryScan:
		return set.RepositoryScan

	case style.ExternalCheck:
		return set.ExternalCheck

	default:
		return nil
	}
}
