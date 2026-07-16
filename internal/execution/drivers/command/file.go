package command

import (
	"errors"
	"fmt"
	"time"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/process"
	"ciphera/tools/internal/style"
)

// runFileCommand runs a file-command tool over its file set. For check execution, the driver
// looks up a FileInterpreter for the tool and converts its raw output into diagnostics; a tool
// without an interpreter is rejected as unsupported. For fix execution (interpreters empty), the
// driver runs the tool and returns empty success on exit 0, or an error otherwise. Fix tools do
// not produce findings to interpret.
func runFileCommand(
	context execution.Context,
	job style.Job,
	interpreters driverkit.FileInterpreters,
	isFix bool,
) (result style.ExecutionResult, err error) {
	fileCommand, found := job.(style.FileCommandExecution)
	if !found {
		return style.ExecutionResult{}, errors.New("file-command driver received empty job")
	}

	files, err := execution.CollectFileSetFiles(context, fileCommand.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return style.ExecutionResult{}, nil
	}

	tool, found := context.Tools[fileCommand.ToolID]
	if !found {
		return style.ExecutionResult{}, errUnknownTool(fileCommand.ToolID)
	}

	arguments := execution.FileCommandArguments(context.RepoRoot, job)
	arguments = append(arguments, files...)
	commandResult, runErr := process.RunCommand(process.CommandRequest{
		Name:             tool.Command,
		Arguments:        arguments,
		Environment:      context.ToolEnvironment,
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

	interpreter, found := interpreters.Lookup(fileCommand.ToolID)
	if !found {
		return style.ExecutionResult{}, fmt.Errorf(
			"no interpreter registered for file-command tool %q",
			fileCommand.ToolID,
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

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
