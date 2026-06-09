package target

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
)

func CheckDrivers(
	commands binding.TargetCommands,
	checks binding.TargetChecks,
) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionTargetCommand: targetCommandDriver(commands),
		style.ExecutionTargetCheck:   targetCheckDriver(checks),
	}
}

func FixDrivers(commands binding.TargetCommands) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionTargetCommand: targetCommandDriver(commands),
	}
}
