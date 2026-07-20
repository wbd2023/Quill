package command

import (
	"context"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// CheckDriver returns the file-command driver for check execution. The driver looks up a
// FileInterpreter for each rule's tool; tools without an interpreter are rejected as
// unsupported rather than silently dumping raw output.
func CheckDriver(interpreters driverkit.FileInterpreters) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		return runFileCommand(ctx, context, job, interpreters, false)
	}
}

// FixDriver returns the file-command driver for fix execution. Fixes never interpret output:
// they either succeed (exit 0, empty result) or fail (non-zero exit, error).
func FixDriver() (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		return runFileCommand(ctx, context, job, driverkit.FileInterpreters{}, true)
	}
}
