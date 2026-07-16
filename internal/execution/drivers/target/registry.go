package target

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
)

// CheckCommandDriver returns the target-command driver for check execution.
func CheckCommandDriver(commands runtimebinding.TargetCommands) (driver execution.Driver) {
	return targetCommandDriver(commands)
}

// CheckCheckDriver returns the target-check driver for check execution.
func CheckCheckDriver(checks runtimebinding.TargetChecks) (driver execution.Driver) {
	return targetCheckDriver(checks)
}

// FixCommandDriver returns the target-command driver for fix execution.
func FixCommandDriver(commands runtimebinding.TargetCommands) (driver execution.Driver) {
	return targetCommandDriver(commands)
}
