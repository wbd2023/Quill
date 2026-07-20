package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
)

func validateStyleGuide(styleGuide policy.StyleGuideConfig) (err error) {
	if isBlank(styleGuide.Path) {
		return fmt.Errorf("style_guide.path must not be empty")
	}

	return nil
}
