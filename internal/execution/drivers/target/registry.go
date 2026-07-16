package target

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
)

// CheckCommandDriver returns the target-command driver for check execution.
func CheckCommandDriver(commands driverkit.TargetCommands) (driver execution.Executor) {
	return targetCommandDriver(commands)
}

// CheckCheckDriver returns the target-check driver for check execution.
func CheckCheckDriver(checks driverkit.TargetChecks) (driver execution.Executor) {
	return targetCheckDriver(checks)
}

// FixCommandDriver returns the target-command driver for fix execution.
func FixCommandDriver(commands driverkit.TargetCommands) (driver execution.Executor) {
	return targetCommandDriver(commands)
}
