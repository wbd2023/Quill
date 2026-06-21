package drivers

import (
	"fmt"

	"ciphera/tools/internal/runner"
	commanddrivers "ciphera/tools/internal/runner/drivers/command"
	profiledrivers "ciphera/tools/internal/runner/drivers/profile"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
	"ciphera/tools/internal/style"
)

// CheckDrivers check drivers.
func CheckDrivers(bindings Bindings) (registry runner.DriverRegistry) {
	return mergeDrivers(
		runner.DriverRegistry{
			style.ExecutionToolchain: runner.ToolchainDriver,
		},
		commanddrivers.CheckDrivers(),
		profiledrivers.CheckDrivers(bindings.projectChecks),
		targetdrivers.CheckDrivers(bindings.targetCommands, bindings.targetChecks),
		scandrivers.CheckDrivers(bindings.repositoryScanners),
	)
}

// FixDrivers fix drivers.
func FixDrivers(bindings Bindings) (registry runner.DriverRegistry) {
	return mergeDrivers(
		commanddrivers.FixDrivers(),
		targetdrivers.FixDrivers(bindings.targetCommands),
	)
}

func mergeDrivers(registries ...runner.DriverRegistry) (merged runner.DriverRegistry) {
	merged = runner.DriverRegistry{}
	for _, registry := range registries {
		for kind, driver := range registry {
			if _, exists := merged[kind]; exists {
				panic(fmt.Sprintf("duplicate driver for execution kind %q", kind))
			}

			merged[kind] = driver
		}
	}

	return merged
}
