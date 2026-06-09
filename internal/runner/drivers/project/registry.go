package project

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
)

func CheckDrivers(checks binding.ProjectChecks) (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		style.ExecutionProject: projectDriver(checks),
	}
}
