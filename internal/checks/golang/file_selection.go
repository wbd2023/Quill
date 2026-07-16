package golang

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func goFilesInDirectories(
	directories []string,
	repository policy.RepositoryConfig,
) (paths []string, err error) {
	return filewalk.CollectFilesInRoots(
		filewalk.WalkConfig{
			ExcludedDirectories: repository.ExcludedDirectories,
			GeneratedMarker:     repository.GeneratedMarker,
		},
		directories,
		".go",
	)
}
