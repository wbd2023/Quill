package builtin

import "ciphera/tools/internal/rules/vocabulary"

func vocabularyPack() (pack Pack) {
	return Pack{
		ID:   PackVocabulary,
		Name: "Vocabulary",
		Config: PackConfig{
			Required: true,
			Validate: vocabulary.ValidatePackConfig,
		},
		Rules: []RuleDefinition{
			scanRule(
				"vocabulary/project-terms",
				"Project vocabulary",
				RuleGroupVocabulary,
				ScannerVocabulary,
			),
		},
	}
}
