package golang

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func goFilesInDirectories(
	directories []string,
	repository policy.RepositoryConfig,
) (paths []string, err error) {
	return filewalk.CollectFilesInRoots(repository, directories, ".go")
}
