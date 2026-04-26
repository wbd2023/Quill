package executors

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

const goLanguage = "go"

func backendCommandExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	detail, found := spec.BackendCommandExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"backend command executor received empty spec",
		)
	}

	switch detail.Action {
	case rulepack.BackendActionGolangci:
		return runGolangci(context, spec)
	case rulepack.BackendActionGoFormat:
		return runGoFormat(context, spec)
	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unknown backend command action %q",
			detail.Action,
		)
	}
}

func backendCheckExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	detail, found := spec.BackendCheckExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("backend check executor received empty spec")
	}

	switch detail.Language {
	case goLanguage:
		return runGoStyleCheck(context, spec)
	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unsupported backend check language %q",
			detail.Language,
		)
	}
}
