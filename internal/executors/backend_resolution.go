package executors

import (
	"fmt"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
)

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

func errEmptyBackendAction(action string) (err error) {
	return fmt.Errorf("%s action received empty spec", action)
}
