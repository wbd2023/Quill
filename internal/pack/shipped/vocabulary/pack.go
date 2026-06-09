package vocabulary

import (
	vocabularyrules "ciphera/tools/internal/checks/vocabulary"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/style"
)

const PackID = "vocabulary"

const ScannerVocabulary = "vocabulary"

const ruleGroupVocabulary style.RuleGroup = "vocabulary_scanners"

// Pack returns the Vocabulary Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Vocabulary",
		Config: pack.Config{
			Required: true,
			Validate: vocabularyrules.ValidatePackConfig,
		},
		Rules: []style.RuleDefinition{
			{
				ID:    "vocabulary/project-terms",
				Name:  "Project vocabulary",
				Group: ruleGroupVocabulary,
				Check: style.ExecutionSpec{
					Kind: style.ExecutionRepositoryScan,
					Detail: style.RepositoryScanExecution{
						Scanner: ScannerVocabulary,
					},
				},
			},
		},
	}
}
