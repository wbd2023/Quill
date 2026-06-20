package styleguide

import (
	"fmt"

	"ciphera/tools/internal/style"
)

// config constants.
const (
	defaultFilename = "STYLE.md"
)

// Config controls STYLE.md parsing.
type Config struct {
	// Filename names the STYLE.md file relative to the repository root.
	// Load requires it; Parse uses it only in diagnostics and defaults to STYLE.md.
	Filename string

	// IDScheme selects the grammar used for requirement IDs.
	IDScheme style.IDScheme
}

func validateIDScheme(scheme style.IDScheme) (err error) {
	if scheme == "" {
		return fmt.Errorf("styleguide requirement id scheme must not be empty")
	}

	if scheme != style.SectionSlug {
		return fmt.Errorf(
			"unsupported styleguide requirement id scheme %q",
			scheme,
		)
	}

	return nil
}
