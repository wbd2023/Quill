package scan

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
)

// CheckDriver returns the repository-scan driver for check execution.
func CheckDriver(scanners runtimebinding.RepositoryScanners) (driver runner.Driver) {
	return repositoryScanDriver(scanners)
}
