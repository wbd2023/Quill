package command

import (
	"errors"
	"fmt"
	"time"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
)

// runFileCommand runs a file-command tool over its file set. For check execution, the driver
// looks up a FileInterpreter for the tool and converts its raw output into diagnostics; a tool
// without an interpreter is rejected as unsupported. For fix execution (interpreters empty), the
// driver runs the tool and returns empty success on exit 0, or an error otherwise. Fix tools do
// not produce findings to interpret.
func runFileCommand(
	context runner.Context,
	spec style.ExecutionSpec,
	interpreters runtimebinding.FileInterpreters,
	isFix bool,
) (result style.ExecutionResult, err error) {
	execution, found := spec.Detail.(style.FileCommandExecution)
	if !found {
		return style.ExecutionResult{}, errors.New("file-command driver received empty spec")
	}

	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return style.ExecutionResult{}, nil
	}

	tool, found := context.Tools[execution.ToolID]
	if !found {
		return style.ExecutionResult{}, errUnknownTool(execution.ToolID)
	}

	arguments := runner.FileCommandArguments(context.RepoRoot, spec)
	arguments = append(arguments, files...)
	commandResult, runErr := runtime.RunCommand(runtime.CommandRequest{
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

	interpreter, found := interpreters.Lookup(execution.ToolID)
	if !found {
		return style.ExecutionResult{}, fmt.Errorf(
			"no interpreter registered for file-command tool %q",
			execution.ToolID,
		)
	}

	if runErr != nil {
		var cmdErr runtime.CommandError
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
