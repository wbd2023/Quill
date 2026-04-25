package rulepack

/* --------------------------------------- Repository Pack -------------------------------------- */

func repositoryPack() (pack Pack) {
	return Pack{
		ID:   PackRepository,
		Name: "Repository",
		Rules: []RuleDefinition{
			repoScanRuleWithConfig(
				"go/architecture-imports",
				"Architecture imports",
				RepositoryScannerArchitecture,
				ConfigRefArchitecture,
			),
			repoScanRuleWithConfig(
				"repo/naming",
				"Naming conventions",
				RepositoryScannerNaming,
				ConfigRefNaming,
			),
		},
	}
}
