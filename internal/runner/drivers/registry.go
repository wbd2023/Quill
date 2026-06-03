package drivers

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

var errViolationsFound = errors.New("violations found")

func CheckDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionToolchain:      runner.ToolchainDriver,
		contract.ExecutionProject:        projectDriver,
		contract.ExecutionFileCommand:    fileCommandDriver,
		contract.ExecutionTargetCommand:  targetCommandDriver,
		contract.ExecutionTargetCheck:    targetCheckDriver,
		contract.ExecutionRepositoryScan: repositoryScanDriver,
	}
}

func FixDrivers() (registry runner.DriverRegistry) {
	return runner.DriverRegistry{
		contract.ExecutionFileCommand:   fileCommandDriver,
		contract.ExecutionTargetCommand: targetCommandDriver,
	}
}
