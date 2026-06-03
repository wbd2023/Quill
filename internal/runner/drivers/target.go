package drivers

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

const goLanguage = "go"

func targetCommandDriver(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.TargetCommandExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"target command driver received empty spec",
		)
	}

	switch execution.Action {
	case builtin.TargetActionGolangci:
		return runGolangci(context, spec)
	case builtin.TargetActionGoFormat:
		return runGoFormat(context, spec)
	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unknown target command action %q",
			execution.Action,
		)
	}
}

func targetCheckDriver(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.TargetCheckExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("target check driver received empty spec")
	}

	switch execution.Language {
	case goLanguage:
		return runGoStyleCheck(context, spec)
	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unsupported target check language %q",
			execution.Language,
		)
	}
}
