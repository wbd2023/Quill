package command

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionFileCommand: fileCommandDriver,
	}
}

func FixDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionFileCommand: fileCommandDriver,
	}
}
