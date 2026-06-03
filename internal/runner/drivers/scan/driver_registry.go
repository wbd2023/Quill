package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionRepositoryScan: repositoryScanDriver,
	}
}
