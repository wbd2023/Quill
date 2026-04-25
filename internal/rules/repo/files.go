package repostyle

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/profile"
)

func ValidateCollectorPolicy(repository profile.RepositoryConfig) (err error) {
	return filewalk.ValidateCollectorPolicy(repository)
}

func CollectFiles(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
	extensions ...string,
) (paths []string, err error) {
	return filewalk.CollectFiles(repoRoot, repository, scope, extensions...)
}

func CollectAllFiles(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (paths []string, err error) {
	return filewalk.CollectAllFiles(repoRoot, repository, scope)
}

func RelativePath(repoRoot string, path string) (relative string) {
	return filewalk.RelativePath(repoRoot, path)
}
