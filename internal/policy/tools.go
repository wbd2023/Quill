package policy

type ToolPin struct {
	ID               string
	Version          string
	TimeoutSeconds   int
	OutputLimitBytes int64
}
