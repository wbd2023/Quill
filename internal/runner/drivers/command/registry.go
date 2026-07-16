package command

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// CheckDriver returns the file-command driver for check execution. The driver looks up a
// FileInterpreter for each rule's tool; tools without an interpreter are rejected as
// unsupported rather than silently dumping raw output.
func CheckDriver(interpreters runtimebinding.FileInterpreters) (driver runner.Driver) {
	return func(
		context runner.Context,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		return runFileCommand(context, job, interpreters, false)
	}
}

// FixDriver returns the file-command driver for fix execution. Fixes never interpret output:
// they either succeed (exit 0, empty result) or fail (non-zero exit, error).
func FixDriver() (driver runner.Driver) {
	return func(
		context runner.Context,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		return runFileCommand(context, job, runtimebinding.FileInterpreters{}, true)
	}
}
