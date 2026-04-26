package policy

type FormattingConfig struct {
	SectionHeaders SectionHeaderConfig
}

type SectionHeaderConfig struct {
	RequiredMinLines  int
	ShortFileMaxLines int
	OveruseCount      int
	GenericNames      []string
	StructuralNames   []string
}
