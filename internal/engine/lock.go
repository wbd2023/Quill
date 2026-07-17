package engine

import (
	"context"

	"ciphera/tools/internal/installer"
	"ciphera/tools/internal/lockfile"
)

// LockResult contains resolved lock file archives.
type LockResult struct {
	Archives []lockfile.Archive
}

// Lock loads the repository profile, resolves every tool's platform archive, and returns
// the entries that make up quill.lock.
func (engine *Engine) Lock(
	operationContext context.Context,
) (result LockResult, operationError error) {
	context, _, err := engine.prepareRunnerContext(operationContext, "")
	if err != nil {
		return LockResult{}, err
	}

	tools := sortedTools(context.Tools)
	archives, err := installer.Resolve(engine.progressWriter, tools)
	if err != nil {
		return LockResult{}, err
	}

	return LockResult{Archives: archives}, nil
}
