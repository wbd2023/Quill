package executors

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	gostyle "ciphera/tools/internal/rules/go"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

/* ---------------------------------------- Go Execution ---------------------------------------- */

func golangciExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	backend, err := goLanguageBackend(context.Policy, spec.Backend)
	if err != nil {
		return "", err
	}

	workdir := languageBackendWorkdir(context.RepoRoot, backend)
	if output, err = runGoFormatChecks(
		workdir,
		context.GoEnvironment,
		context.Policy.Imports.LocalPrefix,
		backend.FormatPaths,
	); err != nil {
		return output, err
	}

	return runtime.RunCommand(
		workdir,
		context.GoEnvironment,
		"golangci-lint",
		"run",
		"./...",
	)
}

func goStyleExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	backend, err := goLanguageBackend(context.Policy, spec.Backend)
	if err != nil {
		return "", err
	}

	if len(backend.StylePaths) == 0 {
		return "", fmt.Errorf("go style backend %q has no paths", spec.Backend)
	}

	directories := make([]string, 0, len(backend.StylePaths))
	for _, current := range backend.StylePaths {
		directories = append(directories, filepath.Join(context.RepoRoot, current))
	}

	return gostyle.CheckDirectories(context.RepoRoot, directories, context.Policy)
}

func goFormatExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	backend, err := goLanguageBackend(context.Policy, spec.Backend)
	if err != nil {
		return "", err
	}

	if len(backend.FormatPaths) == 0 {
		return "", nil
	}

	workdir := languageBackendWorkdir(context.RepoRoot, backend)
	if output, err = runtime.RunCommand(
		workdir,
		context.GoEnvironment,
		"gofmt",
		append([]string{"-w"}, backend.FormatPaths...)...,
	); err != nil {
		return output, err
	}

	return runtime.RunCommand(
		workdir,
		context.GoEnvironment,
		"goimports",
		append(
			[]string{"-w", "-local", context.Policy.Imports.LocalPrefix},
			backend.FormatPaths...,
		)...,
	)
}

/* ------------------------------------- Backend Resolution ------------------------------------- */

func goLanguageBackend(
	policy profile.Profile,
	name string,
) (backend profile.LanguageBackendConfig, err error) {
	backend, found := policy.LanguageBackend(name)
	if !found {
		return profile.LanguageBackendConfig{}, fmt.Errorf("unknown Go backend %q", name)
	}

	if backend.Language != "go" {
		return profile.LanguageBackendConfig{}, fmt.Errorf(
			"backend %q is %q, not go",
			name,
			backend.Language,
		)
	}

	return backend, nil
}

func languageBackendWorkdir(
	repoRoot string,
	backend profile.LanguageBackendConfig,
) (workdir string) {
	if backend.Workdir == "" || backend.Workdir == "." {
		return repoRoot
	}

	return filepath.Join(repoRoot, backend.Workdir)
}

/* ---------------------------------------- Format Checks --------------------------------------- */

func runGoFormatChecks(
	workdir string,
	environment map[string]string,
	localPrefix string,
	paths []string,
) (output string, err error) {
	if len(paths) == 0 {
		return "", nil
	}

	if output, err = runtime.RunCommand(
		workdir,
		environment,
		"gofmt",
		append([]string{"-l"}, paths...)...,
	); err != nil {
		return output, err
	}

	if strings.TrimSpace(output) != "" {
		return "Go files require gofmt formatting:\n" + strings.TrimSpace(output),
			errors.New("gofmt formatting required")
	}

	if output, err = runtime.RunCommand(
		workdir,
		environment,
		"goimports",
		append([]string{"-l", "-local", localPrefix}, paths...)...,
	); err != nil {
		return output, err
	}

	if strings.TrimSpace(output) != "" {
		return "Go files require goimports formatting:\n" + strings.TrimSpace(output),
			errors.New("goimports formatting required")
	}

	return "", nil
}
