package scan

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
)

// CheckDriver returns the repository-scan driver for check execution.
func CheckDriver(scanners driverkit.RepositoryScanners) (driver execution.Executor) {
	return repositoryScanDriver(scanners)
}
