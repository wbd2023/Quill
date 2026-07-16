package drivers

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/command"
	"ciphera/tools/internal/execution/drivers/profile"
	"ciphera/tools/internal/execution/drivers/scan"
	"ciphera/tools/internal/execution/drivers/target"
)

// CheckDrivers returns the complete driver set for check execution.
func CheckDrivers(bindings Bindings) (set execution.DriverSet) {
	return execution.DriverSet{
		Toolchain:      execution.ToolchainDriver,
		Profile:        profile.CheckDriver(bindings.projectChecks),
		FileCommand:    command.CheckDriver(bindings.fileInterpreters),
		TargetCommand:  target.CheckCommandDriver(bindings.targetCommands),
		TargetCheck:    target.CheckCheckDriver(bindings.targetChecks),
		RepositoryScan: scan.CheckDriver(bindings.repositoryScanners),
	}
}

// FixDrivers returns the driver set for fix execution (command and target only).
func FixDrivers(bindings Bindings) (set execution.DriverSet) {
	return execution.DriverSet{
		FileCommand:   command.FixDriver(),
		TargetCommand: target.FixCommandDriver(bindings.targetCommands),
	}
}
