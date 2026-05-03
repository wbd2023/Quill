package policy

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
