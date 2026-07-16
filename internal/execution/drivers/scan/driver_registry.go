package scan

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
)

// CheckDriver returns the repository-scan driver for check execution.
func CheckDriver(scanners runtimebinding.RepositoryScanners) (driver execution.Driver) {
	return repositoryScanDriver(scanners)
}
