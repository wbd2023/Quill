package engine

import (
	"context"

	"github.com/wbd2023/quill/internal/installer"
	"github.com/wbd2023/quill/internal/lockfile"
)

// LockResult describes a completed lock operation.
type LockResult struct {
	// Path is the absolute path of the written quill.lock.
	Path string
	// ArchiveCount is the number of archive-tool entries written.
	ArchiveCount int
}

// Lock loads the repository profile, resolves every tool's platform archive,
// writes quill.lock atomically, and returns the written path and archive count.
func (engine *Engine) Lock(
	operationContext context.Context,
) (result LockResult, operationError error) {
	runContext, _, err := engine.prepareRun(operationContext, "")
	if err != nil {
		return LockResult{}, err
	}

	tools := sortedTools(runContext.Tools)
	archives, err := installer.Resolve(operationContext, engine.progressWriter, tools)
	if err != nil {
		return LockResult{}, err
	}

	return writeResolvedLock(operationContext, engine.repositoryRoot, archives)
}

func writeResolvedLock(
	operationContext context.Context,
	repositoryRoot string,
	archives []lockfile.Archive,
) (result LockResult, operationError error) {
	archiveByID := make(map[string]lockfile.Archive, len(archives))
	for _, archive := range archives {
		archiveByID[archive.Tool] = archive
	}

	if err := operationContext.Err(); err != nil {
		return LockResult{}, err
	}

	path, err := lockfile.Write(operationContext, repositoryRoot, lockfile.Lockfile{Archives: archiveByID})
	if err != nil {
		return LockResult{}, err
	}

	return LockResult{Path: path, ArchiveCount: len(archives)}, nil
}
