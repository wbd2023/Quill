package drivers

import "github.com/wbd2023/quill/internal/execution"

// Checks returns the complete driver set for check execution.
func Checks(bindings Bindings) (set execution.DriverSet) {
	return execution.DriverSet{
		Toolchain:      ToolchainDriver,
		Profile:        profileCheckDriver(bindings.profileChecks),
		FileCommand:    fileCommandCheckDriver(bindings.fileInterpreters),
		TargetCommand:  targetCommandDriver(bindings.targetCommands),
		TargetCheck:    targetCheckDriver(bindings.targetChecks),
		RepositoryScan: repositoryScanDriver(bindings.repositoryScanners),
		ExternalCheck:  externalCheckDriver(),
	}
}

// Fixes returns the driver set for fix execution (command and target only).
func Fixes(bindings Bindings) (set execution.DriverSet) {
	return execution.DriverSet{
		FileCommand:   fileCommandFixDriver(),
		TargetCommand: targetCommandDriver(bindings.targetCommands),
	}
}
