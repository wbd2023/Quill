package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
)

// indexToolIDs builds a set of valid tool IDs for reference checking.
func indexToolIDs(toolIDs []string) (available map[string]bool) {
	available = make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		available[toolID] = true
	}
	return available
}

// validatePins checks that every pinned tool references a real definition and every definition
// has a pin.
func validatePins(config policy.Config, available map[string]bool) (err error) {
	for _, pinnedTool := range config.Tools {
		if !available[pinnedTool.ID] {
			return fmt.Errorf(
				"pinned tool %q does not match an active tool definition",
				pinnedTool.ID,
			)
		}
	}

	for toolID := range available {
		if _, found := config.Tools.Lookup(toolID); !found {
			return fmt.Errorf("active tool %q is missing a pinned tool", toolID)
		}
	}

	return nil
}
