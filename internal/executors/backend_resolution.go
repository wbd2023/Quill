package executors

import (
	"fmt"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

/* ------------------------------------- Backend Resolution ------------------------------------- */

func goLanguageBackends(
	context runner.Context,
	spec contract.ExecutionSpec,
) (backends []policy.LanguageBackendConfig, err error) {
	for _, name := range spec.Backends() {
		backend, err := goLanguageBackend(context.Policy, name)
		if err != nil {
			return nil, err
		}

		if !context.Policy.Repository.ScopesOverlap(context.Scope, backend.Scope) {
			continue
		}

		backends = append(backends, backend)
	}

	return backends, nil
}

func goLanguageBackend(
	config policy.Config,
	name string,
) (backend policy.LanguageBackendConfig, err error) {
	backend, found := config.LanguageBackend(name)
	if !found {
		return policy.LanguageBackendConfig{}, fmt.Errorf("unknown Go backend %q", name)
	}

	if backend.Language != goLanguage {
		return policy.LanguageBackendConfig{}, fmt.Errorf(
			"backend %q is %q, not go",
			name,
			backend.Language,
		)
	}

	return backend, nil
}

func languageBackendWorkdir(
	repoRoot string,
	backend policy.LanguageBackendConfig,
) (workdir string) {
	if backend.Workdir == "" || backend.Workdir == "." {
		return repoRoot
	}

	return filepath.Join(repoRoot, backend.Workdir)
}

/* --------------------------------------- Output Helpers --------------------------------------- */

func appendExecutorOutput(builder *strings.Builder, output string) {
	output = strings.TrimSpace(output)
	if output == "" {
		return
	}

	if builder.Len() > 0 {
		builder.WriteString("\n")
	}

	builder.WriteString(output)
}

/* --------------------------------------- Command Helpers -------------------------------------- */

func runToolByID(
	context runner.Context,
	workdir string,
	toolID string,
	arguments ...string,
) (output string, err error) {
	tool, found := context.Effective.ToolByID(toolID)
	if !found {
		return "", fmt.Errorf("unknown tool %q", toolID)
	}

	capability, found := context.ToolCapabilities[toolID]
	if !found {
		return "", fmt.Errorf("unknown tool capability %q", toolID)
	}

	return runtime.RunToolCommand(workdir, context.GoEnvironment, tool, capability, arguments...)
}

func runCommandOutput(
	workdir string,
	environment map[string]string,
	name string,
	arguments ...string,
) (output string, err error) {
	result, err := runtime.RunCommand(runtime.CommandRequest{
		Directory:   workdir,
		Environment: environment,
		Name:        name,
		Arguments:   append([]string{}, arguments...),
	})
	return runtime.CommandOutput(result, err)
}

func errEmptyBackendAction(action string) (err error) {
	return fmt.Errorf("%s action received empty spec", action)
}
