package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validateTools(tools []policy.PinnedTool) (err error) {
	seen := make(map[string]bool, len(tools))
	for _, tool := range tools {
		if tool.ID == "" {
			return fmt.Errorf("pinned tool has an empty id")
		}

		if seen[tool.ID] {
			return fmt.Errorf("duplicate pinned tool %q", tool.ID)
		}

		seen[tool.ID] = true
		if tool.Version == "" {
			return fmt.Errorf("pinned tool %q must define version", tool.ID)
		}
	}

	return nil
}
