package vocabulary

import (
	"github.com/wbd2023/quill/internal/pack"
	vocabularypolicy "github.com/wbd2023/quill/internal/pack/shipped/vocabulary/policy"
	"github.com/wbd2023/quill/internal/style"
)

// PackID is the canonical identifier for this Pack.
const PackID = "vocabulary"

// Scanner is scanner vocabulary.
const Scanner = "vocabulary"

const ruleGroupVocabulary style.RuleGroup = "vocabulary_scanners"

// Pack returns the Vocabulary Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Vocabulary",
		Policy: pack.Policy{
			Required: true,
			Validate: vocabularypolicy.Validate,
		},
		Rules: []style.RuleDefinition{
			{
				ID:    "vocabulary/project-terms",
				Name:  "Project vocabulary",
				Group: ruleGroupVocabulary,
				Check: style.RepositoryScan{
					Scanner: Scanner,
				},
			},
		},
	}
}
