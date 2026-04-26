package rulepack

func namingPack() (pack Pack) {
	return Pack{
		ID:   PackNaming,
		Name: "Naming",
		Rules: []RuleDefinition{
			scanRuleWithConfig(
				"naming/vocabulary",
				"Naming conventions",
				RuleGroupNaming,
				ScannerNaming,
				ConfigRefNaming,
			),
		},
	}
}
