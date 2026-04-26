package golang

import (
	"path/filepath"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func goFilesInDirectories(
	directories []string,
	repository policy.RepositoryConfig,
) (paths []string, err error) {
	return filewalk.CollectFilesInRoots(repository, directories, ".go")
}

func relativePath(repoRoot string, path string) (relative string) {
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(relative)
}
