package profile

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
)

// CheckDriver returns the profile driver for check execution.
func CheckDriver(checks runtimebinding.ProfileChecks) (driver runner.Driver) {
	return profileDriver(checks)
}
