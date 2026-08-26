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
	ctx context.Context,
) (result LockResult, err error) {
	runContext, _, err := engine.prepareRun(ctx, "")
	if err != nil {
		return LockResult{}, err
	}

	tools := sortedTools(runContext.Tools)
	archives, err := installer.Resolve(ctx, engine.progressWriter, tools)
	if err != nil {
		return LockResult{}, err
	}

	return writeResolvedLock(ctx, engine.root, archives)
}

func writeResolvedLock(
	ctx context.Context,
	root string,
	archives []lockfile.Archive,
) (result LockResult, err error) {
	archiveByID := make(map[string]lockfile.Archive, len(archives))
	for _, archive := range archives {
		archiveByID[archive.Tool] = archive
	}

	if err := ctx.Err(); err != nil {
		return LockResult{}, err
	}

	path, err := lockfile.Write(ctx, root, lockfile.Lockfile{
		Archives: archiveByID,
	})
	if err != nil {
		return LockResult{}, err
	}

	return LockResult{Path: path, ArchiveCount: len(archives)}, nil
}
