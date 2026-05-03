package policy

// VocabularyConfig defines project naming vocabulary rules.
type VocabularyConfig struct {
	Go    GoVocabularyConfig
	Shell ShellVocabularyConfig
}

// GoVocabularyConfig defines vocabulary policy for Go source files.
type GoVocabularyConfig struct {
	ForbiddenTypeSuffixes       []string
	PreferredTypeSuffix         string
	ForbiddenIdentifierSuffixes []string
	PreferredIdentifierSuffix   string
}

// ShellVocabularyConfig defines vocabulary policy for shell scripts.
type ShellVocabularyConfig struct {
	ForbiddenAssignmentNames []string
	PreferredAssignmentName  string
}
