package runner

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func fileSetUsesScopedIncludes(
	fileSet policy.FileSetConfig,
) (usesFilters bool) {
	for _, files := range fileSet.ExplicitFiles {
		if len(files) > 0 {
			return true
		}
	}

	for _, pathPrefixes := range fileSet.PathPrefixes {
		if len(pathPrefixes) > 0 {
			return true
		}
	}

	return false
}

func fileSetIncludeScopes(fileSet policy.FileSetConfig) (scopes []contract.Scope) {
	seen := make(map[contract.Scope]bool)
	for scope, files := range fileSet.ExplicitFiles {
		if len(files) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	for scope, pathPrefixes := range fileSet.PathPrefixes {
		if len(pathPrefixes) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	return scopes
}
