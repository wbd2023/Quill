package target

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
)

// CheckCommandDriver returns the target-command driver for check execution.
func CheckCommandDriver(commands runtimebinding.TargetCommands) (driver runner.Driver) {
	return targetCommandDriver(commands)
}

// CheckCheckDriver returns the target-check driver for check execution.
func CheckCheckDriver(checks runtimebinding.TargetChecks) (driver runner.Driver) {
	return targetCheckDriver(checks)
}

// FixCommandDriver returns the target-command driver for fix execution.
func FixCommandDriver(commands runtimebinding.TargetCommands) (driver runner.Driver) {
	return targetCommandDriver(commands)
}
