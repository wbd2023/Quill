package profile

import "fmt"

// PinnedTools defines the external tools pinned by the profile.
type PinnedTools []PinnedTool

// PinnedTool defines a pinned external tool version and its execution limits.
type PinnedTool struct {
	ID               string
	Version          string
	TimeoutSeconds   int
	OutputLimitBytes int64
}

// Lookup returns the pinned tool with the given ID.
func (t PinnedTools) Lookup(id string) (tool PinnedTool, found bool) {
	for _, candidate := range t {
		if candidate.ID == id {
			return candidate, true
		}
	}

	return PinnedTool{}, false
}

func validateTools(tools []PinnedTool) (err error) {
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
