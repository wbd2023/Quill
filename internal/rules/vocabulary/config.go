package vocabulary

// Config defines project vocabulary policy.
type Config struct {
	Go   GoConfig
	Bash BashConfig
}

// GoConfig defines vocabulary policy for Go source files.
type GoConfig struct {
	ForbiddenTypeSuffixes       []string
	PreferredTypeSuffix         string
	ForbiddenIdentifierSuffixes []string
	PreferredIdentifierSuffix   string
}

// BashConfig defines vocabulary policy for Bash scripts.
type BashConfig struct {
	ForbiddenVariableNames []string
	PreferredVariableName  string
}
