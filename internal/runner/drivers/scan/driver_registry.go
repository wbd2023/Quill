package scan

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
)

func CheckDrivers(scanners binding.RepositoryScanners) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionRepositoryScan: repositoryScanDriver(scanners),
	}
}
