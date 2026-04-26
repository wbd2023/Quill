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
			controlPlaneRule(
				"control-plane/quality-targets",
				"Quality Make targets",
				ControlPlaneCheckQualityTargets,
			),
			controlPlaneRuleWithConfig(
				"control-plane/global-exclusions",
				"Global exclusions",
				ControlPlaneCheckGlobalExclusions,
				ConfigRefRepository,
			),
		},
	}
}
