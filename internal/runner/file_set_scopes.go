package runner

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func fileSetUsesScopedIncludes(
	fileSet policy.FileSetConfig,
) (usesFilters bool) {
	for _, files := range fileSet.Files {
		if len(files) > 0 {
			return true
		}
	}

	for _, prefixes := range fileSet.Prefixes {
		if len(prefixes) > 0 {
			return true
		}
	}

	return false
}

func fileSetIncludeScopes(fileSet policy.FileSetConfig) (scopes []contract.Scope) {
	seen := make(map[contract.Scope]bool)
	for scope, files := range fileSet.Files {
		if len(files) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	for scope, prefixes := range fileSet.Prefixes {
		if len(prefixes) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	return scopes
}
