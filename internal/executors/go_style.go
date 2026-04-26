package executors

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rules/golang"
	"ciphera/tools/internal/runner"
)

func runGoStyleCheck(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.BackendCheckExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("go style check received empty spec")
	}

	backends, err := goLanguageBackends(context, spec)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	diagnostics := make([]contract.Diagnostic, 0)
	var builder strings.Builder
	var joined error
	for _, backend := range backends {
		if len(backend.StylePaths) == 0 {
			joined = errors.Join(
				joined,
				fmt.Errorf("go style backend %q has no paths", backend.Name),
			)
			continue
		}

		roots := make([]string, 0, len(backend.StylePaths))
		for _, stylePath := range backend.StylePaths {
			roots = append(roots, filepath.Join(context.RepoRoot, stylePath))
		}

		styleResult, err := golang.CheckDirectories(
			context.RepoRoot,
			roots,
			context.Policy,
			execution.Check,
		)
		diagnostics = append(diagnostics, styleResult.Diagnostics...)
		appendExecutorOutput(&builder, styleResult.Output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{
		Diagnostics: diagnostics,
		Output:      strings.TrimSpace(builder.String()),
	}, joined
}
