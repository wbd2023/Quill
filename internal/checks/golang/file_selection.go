package golang

import (
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
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
