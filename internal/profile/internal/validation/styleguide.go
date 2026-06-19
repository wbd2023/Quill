package validation

import (
	"fmt"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func validateStyleGuide(styleGuide policy.StyleGuideConfig) (err error) {
	if isBlank(styleGuide.Path) {
		return fmt.Errorf("style_guide.path must not be empty")
	}

	if isBlank(string(styleGuide.IDScheme)) {
		return fmt.Errorf("style_guide.id_scheme must not be empty")
	}

	if styleGuide.IDScheme != style.SectionSlug {
		return fmt.Errorf(
			"unsupported style_guide.id_scheme %q",
			styleGuide.IDScheme,
		)
	}

	return nil
}
