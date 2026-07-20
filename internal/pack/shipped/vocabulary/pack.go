package vocabulary

import (
	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/style"
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
				Check: style.RepositoryScanExecution{
					Scanner: ScannerVocabulary,
				},
			},
		},
	}
}
