package style

// Tool describes a pinned external capability required by one or more Packs.
type Tool struct {
	ID               string
	Name             string
	PinnedVersion    string
	TimeoutSeconds   int
	OutputLimitBytes int64
}
