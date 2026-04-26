package policy

type NamingConfig struct {
	GoTypeSuffixForbidden       []string
	GoTypeSuffixPreferred       string
	GoIdentifierSuffixForbidden []string
	GoIdentifierSuffixPreferred string
	GoParameters                GoParameterConfig
	GoDomainIdentifiers         GoDomainIdentifierConfig
	ShellForbiddenAssignments   []string
	ShellPreferredAssignment    string
}

type GoParameterConfig struct {
	SecretNames           []string
	ConstructorCategories []GoConstructorCategory
}

type GoConstructorCategory struct {
	Name                string
	TypeMarkers         []string
	ExcludedTypeMarkers []string
	ParameterNames      []string
	UsesSecretNames     bool
}

type GoDomainIdentifierConfig map[string][]string
