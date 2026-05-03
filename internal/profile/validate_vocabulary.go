package profile

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"
)

func validateVocabulary(vocabulary policy.VocabularyConfig) (err error) {
	for _, suffix := range vocabulary.Go.ForbiddenTypeSuffixes {
		if strings.TrimSpace(suffix) == "" {
			return fmt.Errorf("vocabulary.go.forbidden_type_suffixes contains an empty suffix")
		}
	}

	for _, suffix := range vocabulary.Go.ForbiddenIdentifierSuffixes {
		if strings.TrimSpace(suffix) == "" {
			return fmt.Errorf(
				"vocabulary.go.forbidden_identifier_suffixes contains an empty suffix",
			)
		}
	}

	for _, name := range vocabulary.Shell.ForbiddenAssignmentNames {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("vocabulary.shell.forbidden_assignment_names contains an empty name")
		}
	}

	return nil
}
