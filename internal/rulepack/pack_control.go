package rulepack

func controlPack() (pack Pack) {
	return Pack{
		ID:    PackControl,
		Name:  "Control plane",
		Tools: coreTools(),
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
			controlPlaneRule(
				"control-plane/enforcement-levels",
				"Enforcement levels",
				ControlPlaneCheckEnforcementLevels,
			),
			controlPlaneRuleWithConfig(
				"control-plane/quality-surface",
				"Quality surface",
				ControlPlaneCheckQualitySurface,
				ConfigReferenceQualitySurface,
			),
			controlPlaneRuleWithConfig(
				"control-plane/global-exclusions",
				"Global exclusions",
				ControlPlaneCheckGlobalExclusions,
				ConfigReferenceRepository,
			),
		},
	}
}
