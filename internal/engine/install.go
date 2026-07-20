package engine

import (
	"context"

	"github.com/wbd2023/Quill/internal/installer"
	"github.com/wbd2023/Quill/internal/lockfile"
	"github.com/wbd2023/Quill/internal/workspace"
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
	context, _, err := engine.prepareRunnerContext(operationContext, "")
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
		sortedTools(context.Tools),
		loaded,
	); err != nil {
		return InstallResult{}, err
	}

	result.Toolchain = engine.inspectTools(
		operationContext,
		context.Tools,
		toolIDs(context.Tools),
		context.ToolEnvironment,
	)
	return result, nil
}
