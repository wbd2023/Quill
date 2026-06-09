package runtime

import "ciphera/tools/internal/toolchain"

func RunToolchainCommand(request toolchain.CommandRequest) (output string, err error) {
	result, err := RunCommand(CommandRequest{
		Directory:   ".",
		Environment: request.Environment,
		Name:        request.Name,
		Arguments:   append([]string{}, request.Arguments...),
	})
	return CommandOutput(result, err)
}
