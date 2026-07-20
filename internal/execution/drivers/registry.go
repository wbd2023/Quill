package drivers

import (
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/command"
	"github.com/wbd2023/Quill/internal/execution/drivers/profile"
	"github.com/wbd2023/Quill/internal/execution/drivers/scan"
	"github.com/wbd2023/Quill/internal/execution/drivers/target"
)

// CheckDrivers returns the complete driver set for check execution.
func CheckDrivers(bindings Bindings) (set execution.DriverSet) {
	return execution.DriverSet{
		Toolchain:      ToolchainDriver,
		Profile:        profile.CheckDriver(bindings.profileChecks),
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
