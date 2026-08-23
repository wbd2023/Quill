package drivers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------ File Command Drivers ------------------------------------ */

// fileCommandCheckDriver returns the file-command driver for check execution. The driver looks up
// a FileInterpreter for each rule's tool; tools without an interpreter are rejected as
// unsupported rather than silently dumping raw output.
func fileCommandCheckDriver(interpreters FileInterpreters) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		_ style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		command, ok := job.(style.FileCommand)
		if !ok {
			return style.ExecutionResult{}, errors.New(
				"file-command driver received unsupported execution job")
		}
		return runFileCommand(ctx, context, command, interpreters, false)
	}
}

// fileCommandFixDriver returns the file-command driver for fix execution. Fixes never interpret
// output: they either succeed (exit 0, empty result) or fail (non-zero exit, error).
func fileCommandFixDriver() (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		_ style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		command, ok := job.(style.FileCommand)
		if !ok {
			return style.ExecutionResult{}, errors.New(
				"file-command driver received unsupported execution job")
		}
		return runFileCommand(ctx, context, command, FileInterpreters{}, true)
	}
}

// runFileCommand runs a file-command tool over its file set. For check execution, the driver
// looks up a FileInterpreter for the tool and converts its raw output into diagnostics; a tool
// without an interpreter is rejected as unsupported. For fix execution (interpreters empty), the
// driver runs the tool and returns empty success on exit 0, or an error otherwise. Fix tools do
// not produce findings to interpret.
func runFileCommand(
	ctx context.Context,
	context execution.RunContext,
	command style.FileCommand,
	interpreters FileInterpreters,
	isFix bool,
) (result style.ExecutionResult, err error) {
	files, err := execution.CollectFileSetFiles(context, command.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return style.ExecutionResult{}, nil
	}

	tool, found := context.Tools[command.ToolID]
	if !found {
		return style.ExecutionResult{}, errUnknownTool(command.ToolID)
	}

	arguments := execution.FileCommandArguments(context.RepoRoot, command, files)
	commandResult, runErr := process.RunCommand(ctx, process.CommandRequest{
		Name:             tool.Command,
		Arguments:        arguments,
		Environment:      process.EnvironmentInherit,
		Variables:        context.ToolEnvironment,
		Directory:        context.RepoRoot,
		Timeout:          time.Duration(tool.TimeoutSeconds) * time.Second,
		OutputLimitBytes: tool.OutputLimitBytes,
	})

	result = style.ExecutionResult{
		ExitCode:  commandResult.ExitCode,
		TimedOut:  commandResult.TimedOut,
		Truncated: commandResult.Truncated,
	}

	if isFix {
		result.Output = commandResult.Output
		return result, runErr
	}

	interpreter, found := interpreters.Lookup(command.ToolID)
	if !found {
		return style.ExecutionResult{}, fmt.Errorf(
			"no interpreter registered for file-command tool %q",
			command.ToolID,
		)
	}

	if runErr != nil {
		var cmdErr process.CommandError
		if !errors.As(runErr, &cmdErr) {
			return result, runErr
		}
	}

	diagnostics, interpErr := interpreter(commandResult)
	result.Diagnostics = diagnostics
	if interpErr != nil {
		return result, interpErr
	}

	return result, nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
