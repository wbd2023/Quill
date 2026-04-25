package rulepack

import "ciphera/tools/internal/contract"

/* ------------------------------------------- Go Pack ------------------------------------------ */

func goPack() (pack Pack) {
	return Pack{
		ID:   PackGo,
		Name: "Go",
		Tools: selectTools(
			contract.ToolGo,
			contract.ToolGoimports,
			contract.ToolGolangciLint,
		),
		Rules: []RuleDefinition{
			golangciRule(
				"go/lint-app",
				"golangci-lint (app)",
				"go_app",
			),
			goStyleRule(
				"go/style-app",
				"Go style (app)",
				"go_app",
			),
			golangciRule(
				"go/lint-tools",
				"golangci-lint (style platform)",
				"go_tools",
			),
			goStyleRule(
				"go/style-tools",
				"Go style (style platform)",
				"go_tools",
			),
			repoScanRule(
				"go/guard-clause-spacing",
				"Guard-clause spacing (Go)",
				RepositoryScannerGuardClauseSpacing,
			),
			repoScanRule(
				"go/switch-case-spacing",
				"Switch-case spacing (Go)",
				RepositoryScannerSwitchCaseSpacing,
			),
		},
	}
}
