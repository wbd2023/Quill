package profile

import "fmt"

func validateStyleGuide(styleGuide StyleGuideConfig) (err error) {
	if isBlank(styleGuide.Path) {
		return fmt.Errorf("style_guide.path must not be empty")
	}

	if err = validateRepoPath(styleGuide.Path); err != nil {
		return fmt.Errorf("style_guide.path: %w", err)
	}

	return nil
}
