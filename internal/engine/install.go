package engine

import (
	"context"

	"github.com/wbd2023/quill/internal/installer"
	"github.com/wbd2023/quill/internal/lockfile"
	"github.com/wbd2023/quill/internal/workspace"
)

// InstallResult contains post-install tool inspection.
type InstallResult struct {
	Toolchain ToolchainInspection
}

// Install loads the repository and lock file, installs configured tools, and inspects the
// resulting toolchain.
func (engine *Engine) Install(
	operationContext context.Context,
) (result InstallResult, operationError error) {
	runContext, _, err := engine.prepareRun(operationContext, "")
	if err != nil {
		return InstallResult{}, err
	}

	layout := workspace.NewLayout(engine.repositoryRoot)
	loaded, err := lockfile.Load(engine.repositoryRoot)
	if err != nil {
		return InstallResult{}, err
	}

	if err = installer.Install(
		operationContext,
		layout,
		engine.progressWriter,
		sortedTools(runContext.Tools),
		loaded,
	); err != nil {
		return InstallResult{}, err
	}

	result.Toolchain, err = engine.inspectTools(
		operationContext,
		runContext.Tools,
		toolIDs(runContext.Tools),
		runContext.ToolEnvironment,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}
