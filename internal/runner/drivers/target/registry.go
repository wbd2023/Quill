package target

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionTargetCommand: targetCommandDriver,
		contract.ExecutionTargetCheck:   targetCheckDriver,
	}
}

func FixDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionTargetCommand: targetCommandDriver,
	}
}
