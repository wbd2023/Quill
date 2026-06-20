package target

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckDrivers check drivers.
func CheckDrivers(
	commands runtimebinding.TargetCommands,
	checks runtimebinding.TargetChecks,
) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionTargetCommand: targetCommandDriver(commands),
		style.ExecutionTargetCheck:   targetCheckDriver(checks),
	}
}

// FixDrivers fix drivers.
func FixDrivers(commands runtimebinding.TargetCommands) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionTargetCommand: targetCommandDriver(commands),
	}
}
