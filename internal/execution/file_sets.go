package execution

import (
	"errors"

	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/policy"
)

// CollectFileSetFiles collect file set files.
func CollectFileSetFiles(context RunContext, name string) (files []string, err error) {
	fileSet, found := context.Profile.FileSets.Lookup(name)
	if !found {
		return nil, errUnknownFileSet(name)
	}

	candidates, err := collectFileSetCandidates(context, fileSet)
	if err != nil {
		return nil, err
	}

	for _, path := range candidates {
		if !fileSetCoversPath(context, fileSet, path) {
			continue
		}

		if fileSetExcludesPath(fileSet, path) {
			continue
		}

		if filewalk.IsBinaryFile(path) {
			continue
		}

		files = append(files, path)
	}

	return files, nil
}

// CollectScopeFiles returns every non-binary file beneath the current scope roots, with no file-set
// filtering. External Pack checks that declare no file set receive this candidate list so the Pack
// can scope its own analysis without re-walking the repository.
func CollectScopeFiles(context RunContext) (files []string, err error) {
	candidates, err := collectFileSetCandidates(context, policy.FileSetConfig{})
	if err != nil {
		return nil, err
	}

	for _, path := range candidates {
		if filewalk.IsBinaryFile(path) {
			continue
		}
		files = append(files, path)
	}

	return files, nil
}

func errUnknownFileSet(name string) (err error) {
	return errors.New("unknown file set " + name)
}
