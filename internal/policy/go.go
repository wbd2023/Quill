package policy

// GoConfig defines Go-specific style policy.
type GoConfig struct {
	LocalImportPrefixes          []string
	Parameters                   GoParameterConfig
	DomainIdentifierConstructors GoDomainIdentifierConstructors
	Architecture                 GoArchitectureConfig
}

// GoParameterConfig defines Go parameter naming and ordering policy.
type GoParameterConfig struct {
	SecretNames      []string
	ConstructorOrder []GoParameterGroup
}

// GoParameterGroup describes a named Go parameter group.
type GoParameterGroup struct {
	Name                string
	TypeMarkers         []string
	ExcludedTypeMarkers []string
	ParameterNames      []string
	MatchesSecretNames  bool
}

// GoDomainIdentifierConstructors maps domain identifier types to approved constructors.
type GoDomainIdentifierConstructors map[string][]string

// GoArchitectureConfig defines package import layers enforced by Go architecture rules.
type GoArchitectureConfig struct {
	Layers []GoArchitectureLayer
}

// GoArchitectureLayer describes one named Go package layer and its allowed dependency layers.
type GoArchitectureLayer struct {
	Name          string
	PackageRoots  []string
	AllowedLayers []string
}
