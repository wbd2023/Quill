package builtin

import "ciphera/tools/internal/rules/project"

func projectPack() (pack Pack) {
	return Pack{
		ID:    PackProject,
		Name:  "Project",
		Tools: coreTools(),
		Config: PackConfig{
			Required: true,
			Validate: project.ValidatePackConfig,
		},
		Rules: []RuleDefinition{
			toolchainRule(
				"toolchain/check-versions",
				"Pinned toolchain versions",
				ToolGo,
				ToolGoimports,
				ToolMisspell,
				ToolGolangciLint,
				ToolShfmt,
				ToolShellcheck,
				ToolMarkdownlint,
			),
			projectRule(
				"project/enforcement-levels",
				"Enforcement levels",
				ProjectCheckEnforcementLevels,
			),
			projectRule(
				"project/quality-commands",
				"Quality commands",
				ProjectCheckCommands,
			),
			projectRule(
				"project/excluded-directories",
				"Excluded directories",
				ProjectCheckExcludedDirectories,
			),
		},
	}
}
