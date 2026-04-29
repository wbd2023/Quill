package styleguide

import (
	"fmt"

	"ciphera/tools/internal/requirementid"
)

const (
	defaultFilename = "STYLE.md"
)

// Config controls STYLE.md parsing.
type Config struct {
	// Filename names the STYLE.md file relative to the repository root.
	// Load requires it; Parse uses it only in diagnostics and defaults to STYLE.md.
	Filename string

	// RequirementIDScheme selects the grammar used for requirement IDs.
	RequirementIDScheme requirementid.Scheme
}

func validateRequirementIDScheme(scheme requirementid.Scheme) (err error) {
	if scheme == "" {
		return fmt.Errorf("styleguide requirement id scheme must not be empty")
	}

	if scheme != requirementid.SectionSlug {
		return fmt.Errorf(
			"unsupported styleguide requirement id scheme %q",
			scheme,
		)
	}

	return nil
}
