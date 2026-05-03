package policy

// FormattingConfig defines formatting policy.
type FormattingConfig struct {
	SectionHeaders SectionHeaderConfig
}

// SectionHeaderConfig defines when block section headers are expected or overused.
type SectionHeaderConfig struct {
	RequiredMinLines  int
	ShortFileMaxLines int
	OveruseThreshold  int
	GenericNames      []string
	StructuralNames   []string
}
