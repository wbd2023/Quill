package drivers

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
	commanddrivers "ciphera/tools/internal/runner/drivers/command"
	projectdrivers "ciphera/tools/internal/runner/drivers/project"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
)

func CheckDrivers() (registry runner.DriverRegistry) {
	return mergeDrivers(
		runner.DriverRegistry{
			contract.ExecutionToolchain: runner.ToolchainDriver,
		},
		commanddrivers.CheckDrivers(),
		projectdrivers.CheckDrivers(),
		targetdrivers.CheckDrivers(),
		scandrivers.CheckDrivers(),
	)
}

func FixDrivers() (registry runner.DriverRegistry) {
	return mergeDrivers(
		commanddrivers.FixDrivers(),
		targetdrivers.FixDrivers(),
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
