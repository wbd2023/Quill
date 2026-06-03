package command

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionFileCommand: fileCommandDriver,
	}
}

func FixDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionFileCommand: fileCommandDriver,
	}
}
