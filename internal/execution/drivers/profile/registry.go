package profile

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
)

// CheckDriver returns the profile driver for check execution.
func CheckDriver(checks runtimebinding.ProfileChecks) (driver execution.Driver) {
	return profileDriver(checks)
}
