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
	ctx context.Context,
) (result InstallResult, err error) {
	runContext, _, err := engine.prepareRun(ctx, "")
	if err != nil {
		return InstallResult{}, err
	}

	layout := workspace.NewLayout(engine.root)
	loaded, err := lockfile.Load(engine.root)
	if err != nil {
		return InstallResult{}, err
	}

	if err = installer.Install(
		ctx,
		layout,
		engine.progressWriter,
		sortedTools(runContext.Tools),
		loaded,
	); err != nil {
		return InstallResult{}, err
	}

	result.Toolchain, err = engine.inspectTools(
		ctx,
		runContext.Tools,
		toolIDs(runContext.Tools),
		runContext.ToolEnvironment,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}
