package engine

import (
	"context"
	"slices"

	"ciphera/tools/internal/coverage"
	"ciphera/tools/internal/installer"
	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/styleguide"
	"ciphera/tools/internal/toolchain"
	"ciphera/tools/internal/workspace"
)

/* ------------------------------------- Coverage Reporting ------------------------------------- */

// Coverage loads STYLE.md and the compiled effective profile, then builds requirement coverage.
//
// Coverage intentionally does not construct a runner context or inspect tools. It uses a separate
// internal pipeline that loads only the profile and style guide.
func (engine *Engine) Coverage(
	operationContext context.Context,
) (coverageReport coverage.Report, operationError error) {
	compiled, err := engine.loadCompiledProfile(operationContext)
	if err != nil {
		return coverage.Report{}, err
	}

	document, err := styleguide.Load(engine.repositoryRoot, styleguide.Config{
		Filename: compiled.profile.Profile.StyleGuide.Path,
	})
	if err != nil {
		return coverage.Report{}, err
	}

	return coverage.Build(document, compiled.profile.Effective.Rules), nil
}

/* -------------------------------------- Tool Installation ------------------------------------- */

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

	tools := sortedTools(context.Tools)

	if err = installer.Install(layout, engine.progressWriter, tools, loaded); err != nil {
		return InstallResult{}, err
	}

	toolIDs := make([]string, 0, len(context.Tools))
	for toolID := range context.Tools {
		toolIDs = append(toolIDs, toolID)
	}

	result.Toolchain = engine.inspectTools(context.Tools, toolIDs, context.ToolEnvironment)
	return result, nil
}

func sortedTools(tools map[string]toolchain.Tool) (sorted []toolchain.Tool) {
	toolIDs := make([]string, 0, len(tools))
	for toolID := range tools {
		toolIDs = append(toolIDs, toolID)
	}

	slices.Sort(toolIDs)
	for _, toolID := range toolIDs {
		sorted = append(sorted, tools[toolID])
	}
	return sorted
}
