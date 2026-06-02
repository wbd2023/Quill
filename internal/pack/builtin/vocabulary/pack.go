package vocabulary

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack"
	vocabularyrules "ciphera/tools/internal/rules/vocabulary"
)

const PackID = "vocabulary"

const ScannerVocabulary = "vocabulary"

const ruleGroupVocabulary contract.RuleGroup = "vocabulary_scanners"

// Pack returns the Vocabulary Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Vocabulary",
		Config: pack.Config{
			Required: true,
			Validate: vocabularyrules.ValidatePackConfig,
		},
		Rules: []contract.RuleDefinition{
			{
				ID:    "vocabulary/project-terms",
				Name:  "Project vocabulary",
				Group: ruleGroupVocabulary,
				Check: contract.ExecutionSpec{
					Kind: contract.ExecutorRepositoryScan,
					Detail: contract.RepositoryScanExecution{
						Scanner: ScannerVocabulary,
					},
				},
			},
		},
	}
}
