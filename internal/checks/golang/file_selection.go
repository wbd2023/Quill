package golang

import (
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/profile"
)

func goFilesInDirectories(
	directories []string,
	repository profile.RepositoryConfig,
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
