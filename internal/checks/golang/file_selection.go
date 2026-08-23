package golang

import (
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/policy"
)

func goFilesInDirectories(
	directories []string,
	repository policy.RepositoryConfig,
) (paths []string, err error) {
	return filewalk.CollectFiles(
		directories,
		filewalk.WalkConfig{
			ExcludedDirectories: repository.ExcludedDirectories,
			GeneratedMarker:     repository.GeneratedMarker,
		},
		".go",
	)
}
