package gopolicy

// Config defines Go rule policy.
type Config struct {
	LocalImportPrefixes []string
	Parameters          ParameterConfig
	Constructors        ConstructorConfig
	DomainValues        DomainValueConfig
	Architecture        ArchitectureConfig
}
