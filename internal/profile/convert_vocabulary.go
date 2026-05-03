package profile

import "ciphera/tools/internal/policy"

func vocabularyFromSchema(schema schemaVocabularyConfig) (vocabulary policy.VocabularyConfig) {
	return policy.VocabularyConfig{
		Go: policy.GoVocabularyConfig{
			ForbiddenTypeSuffixes: append([]string{}, schema.Go.ForbiddenTypeSuffixes...),
			PreferredTypeSuffix:   schema.Go.PreferredTypeSuffix,
			ForbiddenIdentifierSuffixes: append(
				[]string{},
				schema.Go.ForbiddenIdentifierSuffixes...,
			),
			PreferredIdentifierSuffix: schema.Go.PreferredIdentifierSuffix,
		},
		Shell: policy.ShellVocabularyConfig{
			ForbiddenAssignmentNames: append(
				[]string{},
				schema.Shell.ForbiddenAssignmentNames...,
			),
			PreferredAssignmentName: schema.Shell.PreferredAssignmentName,
		},
	}
}

func vocabularyToSchema(vocabulary policy.VocabularyConfig) (schema schemaVocabularyConfig) {
	return schemaVocabularyConfig{
		Go: schemaGoVocabularyConfig{
			ForbiddenTypeSuffixes: append([]string{}, vocabulary.Go.ForbiddenTypeSuffixes...),
			PreferredTypeSuffix:   vocabulary.Go.PreferredTypeSuffix,
			ForbiddenIdentifierSuffixes: append(
				[]string{},
				vocabulary.Go.ForbiddenIdentifierSuffixes...,
			),
			PreferredIdentifierSuffix: vocabulary.Go.PreferredIdentifierSuffix,
		},
		Shell: schemaShellVocabularyConfig{
			ForbiddenAssignmentNames: append(
				[]string{},
				vocabulary.Shell.ForbiddenAssignmentNames...,
			),
			PreferredAssignmentName: vocabulary.Shell.PreferredAssignmentName,
		},
	}
}
