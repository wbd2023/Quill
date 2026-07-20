package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
)

func validateTools(tools []policy.PinnedTool) (err error) {
	seen := make(map[string]bool, len(tools))
	for _, tool := range tools {
		if isBlank(tool.ID) {
			return fmt.Errorf("pinned tool has an empty id")
		}

		if seen[tool.ID] {
			return fmt.Errorf("duplicate pinned tool %q", tool.ID)
		}

		seen[tool.ID] = true
		if isBlank(tool.Version) {
			return fmt.Errorf("pinned tool %q must define version", tool.ID)
		}

		if tool.TimeoutSeconds < 0 {
			return fmt.Errorf("pinned tool %q timeout_seconds must not be negative", tool.ID)
		}

		if tool.OutputLimitBytes < 0 {
			return fmt.Errorf("pinned tool %q output_limit_bytes must not be negative", tool.ID)
		}
	}

	return nil
}
