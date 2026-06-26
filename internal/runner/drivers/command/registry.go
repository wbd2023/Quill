package command

import "ciphera/tools/internal/runner"

// CheckDriver returns the file-command driver for check execution.
func CheckDriver() (driver runner.Driver) { return fileCommandDriver }

// FixDriver returns the file-command driver for fix execution.
func FixDriver() (driver runner.Driver) { return fileCommandDriver }
