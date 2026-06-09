package text

// Config defines text rule policy.
type Config struct {
	SectionHeaders SectionHeaderConfig
}

// SectionHeaderConfig defines when block section headers are expected or overused.
type SectionHeaderConfig struct {
	LargeMinLines   int
	ShortMaxLines   int
	MaxHeaderCount  int
	GenericNames    []string
	StructuralNames []string
}
