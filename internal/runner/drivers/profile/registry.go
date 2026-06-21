package project

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckDrivers check drivers.
func CheckDrivers(checks runtimebinding.ProfileChecks) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionProfile: projectDriver(checks),
	}
}
