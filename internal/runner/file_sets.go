package runner

import (
	"errors"

	"ciphera/tools/internal/filewalk"
)

func CollectFileSetFiles(context Context, name string) (files []string, err error) {
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

func errUnknownFileSet(name string) (err error) {
	return errors.New("unknown file set " + name)
}
