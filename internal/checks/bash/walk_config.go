package bash

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func walkConfig(repository policy.RepositoryConfig) (config filewalk.WalkConfig) {
	return filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}
}
