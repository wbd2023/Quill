package command

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

// CheckDrivers check drivers.
func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionFileCommand: fileCommandDriver,
	}
}

// FixDrivers fix drivers.
func FixDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionFileCommand: fileCommandDriver,
	}
}
