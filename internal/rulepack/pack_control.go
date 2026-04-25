package rulepack

import "ciphera/tools/internal/contract"

/* ---------------------------------------- Control Pack ---------------------------------------- */

func controlPack() (pack Pack) {
	return Pack{
		ID:    PackControl,
		Name:  "Control plane",
		Tools: coreTools(),
		Rules: []RuleDefinition{
			toolchainRule(
				"toolchain/check-versions",
				"Pinned toolchain versions",
				contract.ToolGo,
				contract.ToolGoimports,
				contract.ToolMisspell,
				contract.ToolGolangciLint,
				contract.ToolShfmt,
				contract.ToolShellcheck,
				contract.ToolMarkdownlint,
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
