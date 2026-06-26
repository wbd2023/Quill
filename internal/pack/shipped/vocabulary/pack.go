package vocabulary

import (
	"ciphera/tools/internal/checks/vocabularypolicy"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/style"
)

// PackID is the canonical identifier for this Pack.
const PackID = "vocabulary"

// ScannerVocabulary is scanner vocabulary.
const ScannerVocabulary = "vocabulary"

const ruleGroupVocabulary style.RuleGroup = "vocabulary_scanners"

// Pack returns the Vocabulary Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Vocabulary",
		Config: pack.Config{
			Required: true,
			Validate: vocabularypolicy.ValidatePackConfig,
		},
		Rules: []style.RuleDefinition{
			{
				ID:    "vocabulary/project-terms",
				Name:  "Project vocabulary",
				Group: ruleGroupVocabulary,
				Check: style.ExecutionSpec{
					Detail: style.RepositoryScanExecution{
						Scanner: ScannerVocabulary,
					},
				},
			},
		},
	}
}
