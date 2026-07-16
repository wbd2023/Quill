package profile

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
)

// CheckDriver returns the profile driver for check execution.
func CheckDriver(checks driverkit.ProfileChecks) (driver execution.Executor) {
	return profileDriver(checks)
}
