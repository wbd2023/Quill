package drivers

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/command"
	"ciphera/tools/internal/runner/drivers/profile"
	"ciphera/tools/internal/runner/drivers/scan"
	"ciphera/tools/internal/runner/drivers/target"
)

// CheckDrivers returns the complete driver set for check execution.
func CheckDrivers(bindings Bindings) (set runner.DriverSet) {
	return runner.DriverSet{
		Toolchain:      runner.ToolchainDriver,
		Profile:        profile.CheckDriver(bindings.projectChecks),
		FileCommand:    command.CheckDriver(),
		TargetCommand:  target.CheckCommandDriver(bindings.targetCommands),
		TargetCheck:    target.CheckCheckDriver(bindings.targetChecks),
		RepositoryScan: scan.CheckDriver(bindings.repositoryScanners),
	}
}

// FixDrivers returns the driver set for fix execution (command and target only).
func FixDrivers(bindings Bindings) (set runner.DriverSet) {
	return runner.DriverSet{
		FileCommand:   command.FixDriver(),
		TargetCommand: target.FixCommandDriver(bindings.targetCommands),
	}
}
