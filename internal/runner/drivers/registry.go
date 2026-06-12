package drivers

import (
	"fmt"

	"ciphera/tools/internal/runner"
	commanddrivers "ciphera/tools/internal/runner/drivers/command"
	projectdrivers "ciphera/tools/internal/runner/drivers/project"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
	"ciphera/tools/internal/style"
)

func CheckDrivers(bindings Bindings) (registry runner.DriverRegistry) {
	return mergeDrivers(
		runner.DriverRegistry{
			style.ExecutionToolchain: runner.ToolchainDriver,
		},
		commanddrivers.CheckDrivers(),
		projectdrivers.CheckDrivers(bindings.projectChecks),
		targetdrivers.CheckDrivers(bindings.targetCommands, bindings.targetChecks),
		scandrivers.CheckDrivers(bindings.repositoryScanners),
	)
}

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
