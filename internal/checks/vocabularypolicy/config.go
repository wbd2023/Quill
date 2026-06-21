package vocabularypolicy

// Config defines Vocabulary Pack Policy.
type Config struct {
	Go   GoConfig
	Bash BashConfig
}

// GoConfig defines Vocabulary Pack Policy for Go source files. TypeSuffixes and IdentifierSuffixes
// map a preferred form to the list of. forbidden shorthands that map onto it.
type GoConfig struct {
	TypeSuffixes       map[string][]string
	IdentifierSuffixes map[string][]string
}

// BashConfig defines Vocabulary Pack Policy for Bash scripts. VariableNames maps a preferred name
// to the list of forbidden names that map. onto it.
type BashConfig struct {
	VariableNames map[string][]string
}
