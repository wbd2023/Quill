package profile

type schemaVocabularyConfig struct {
	Go    schemaGoVocabularyConfig    `toml:"go"`
	Shell schemaShellVocabularyConfig `toml:"shell"`
}

type schemaGoVocabularyConfig struct {
	ForbiddenTypeSuffixes       []string `toml:"forbidden_type_suffixes"`
	PreferredTypeSuffix         string   `toml:"preferred_type_suffix"`
	ForbiddenIdentifierSuffixes []string `toml:"forbidden_identifier_suffixes"`
	PreferredIdentifierSuffix   string   `toml:"preferred_identifier_suffix"`
}

type schemaShellVocabularyConfig struct {
	ForbiddenAssignmentNames []string `toml:"forbidden_assignment_names"`
	PreferredAssignmentName  string   `toml:"preferred_assignment_name"`
}
