package style

type Tool struct {
	ID               string
	Name             string
	PinnedVersion    string
	TimeoutSeconds   int
	OutputLimitBytes int64
}
