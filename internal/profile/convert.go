package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func configFromSchema(schema schemaConfig) (config policy.Config) {
	return policy.Config{
		SchemaVersion:  schema.SchemaVersion,
		Repository:     repositoryFromSchema(schema.Repository),
		StyleGuide:     styleGuideFromSchema(schema.StyleGuide),
		Paths:          cloneStringSlices(policy.PathClasses(schema.Paths)),
		FileSets:       fileSetsFromSchema(schema.FileSets),
		Language:       languageFromSchema(schema.Language),
		Go:             goFromSchema(schema.Go),
		Tools:          toolsFromSchema(schema.Tools),
		Formatting:     formattingFromSchema(schema.Formatting),
		Vocabulary:     vocabularyFromSchema(schema.Vocabulary),
		QualitySurface: qualitySurfaceFromSchema(schema.QualitySurface),
		RulePacks: policy.RulePackConfig{
			Enabled: append([]string{}, schema.RulePacks.Enabled...),
		},
		Rules: rulesFromSchema(schema.Rules),
	}
}

func schemaFromConfig(config policy.Config) (schema schemaConfig) {
	return schemaConfig{
		SchemaVersion: config.SchemaVersion,
		Repository:    repositoryToSchema(config.Repository),
		StyleGuide: schemaStyleGuideConfig{
			Path:                config.StyleGuide.Path,
			RequirementIDScheme: string(config.StyleGuide.RequirementIDScheme),
		},
		Paths:          cloneStringSlices(config.Paths),
		FileSets:       fileSetsToSchema(config.FileSets),
		Language:       languageToSchema(config.Language),
		Go:             goToSchema(config.Go),
		Tools:          toolsToSchema(config.Tools),
		Formatting:     formattingToSchema(config.Formatting),
		Vocabulary:     vocabularyToSchema(config.Vocabulary),
		QualitySurface: qualitySurfaceToSchema(config.QualitySurface),
		RulePacks: schemaRulePackConfig{
			Enabled: append([]string{}, config.RulePacks.Enabled...),
		},
		Rules: rulesToSchema(config.Rules),
	}
}

func scopeMapFromSchema(source map[string][]string) (target map[contract.Scope][]string) {
	if source == nil {
		return nil
	}

	target = make(map[contract.Scope][]string, len(source))
	for scope, values := range source {
		target[contract.Scope(scope)] = append([]string{}, values...)
	}

	return target
}

func scopeMapToSchema(source map[contract.Scope][]string) (target map[string][]string) {
	if source == nil {
		return nil
	}

	target = make(map[string][]string, len(source))
	for scope, values := range source {
		target[string(scope)] = append([]string{}, values...)
	}

	return target
}

func cloneStringSlices[M ~map[string][]string](source M) (target M) {
	if source == nil {
		return nil
	}

	target = make(M, len(source))
	for key, values := range source {
		target[key] = append([]string{}, values...)
	}

	return target
}
