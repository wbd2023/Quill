package scan

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckDrivers check drivers.
func CheckDrivers(scanners runtimebinding.RepositoryScanners) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionRepositoryScan: repositoryScanDriver(scanners),
	}
}
