package policy

// Config defines Vocabulary Pack Policy.
type Config struct {
	Go   GoConfig
	Bash BashConfig
}

// GoConfig defines Vocabulary Pack Policy for Go source files.
type GoConfig struct {
	ForbiddenTypeSuffixes       []string
	PreferredTypeSuffix         string
	ForbiddenIdentifierSuffixes []string
	PreferredIdentifierSuffix   string
}

// BashConfig defines Vocabulary Pack Policy for Bash scripts.
type BashConfig struct {
	ForbiddenVariableNames []string
	PreferredVariableName  string
}
