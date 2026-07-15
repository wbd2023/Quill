package command

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
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
	execution, found := spec.FileCommandExecution()
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
		Directory:        context.RepoRoot,
		Environment:      context.ToolEnvironment,
		Name:             tool.Command,
		Arguments:        arguments,
		TimeoutSeconds:   tool.TimeoutSeconds,
		OutputLimitBytes: tool.OutputLimitBytes,
	})

	styleCommand := commandrun.BuildStyleResult(commandResult)

	if isFix {
		return style.ExecutionResult{Command: styleCommand}, runErr
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
			return style.ExecutionResult{Command: styleCommand}, runErr
		}
	}

	diagnostics, interpErr := interpreter(commandResult)
	if interpErr != nil {
		return style.ExecutionResult{Diagnostics: diagnostics, Command: styleCommand}, interpErr
	}

	return style.ExecutionResult{Diagnostics: diagnostics, Command: styleCommand}, nil
}

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
