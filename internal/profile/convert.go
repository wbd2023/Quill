package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Conversion ----------------------------------------- */

func policyFromSchema(schema schemaConfig) (config policy.Config) {
	return policy.Config{
		SchemaVersion: schema.SchemaVersion,
		RulePacks: policy.RulePackConfig{
			Enabled: append([]string{}, schema.RulePacks.Enabled...),
		},
		Repository:   repositoryFromSchema(schema.Repository),
		StyleGuide:   styleGuideFromSchema(schema.StyleGuide),
		Formatting:   formattingFromSchema(schema.Formatting),
		Imports:      policy.ImportsConfig{LocalPrefix: schema.Imports.LocalPrefix},
		Paths:        policy.PathClassSet{Classes: cloneStringMap(schema.Paths)},
		FileSets:     fileSetsFromSchema(schema.FileSets),
		Language:     languageFromSchema(schema.Language),
		Tools:        toolsFromSchema(schema.Tools),
		Naming:       namingFromSchema(schema.Naming),
		ControlPlane: controlPlaneFromSchema(schema.ControlPlane),
		Architecture: architectureFromSchema(schema.Architecture),
		Rules:        rulesFromSchema(schema.Rules),
	}
}

func schemaFromPolicy(config policy.Config) (schema schemaConfig) {
	return schemaConfig{
		SchemaVersion: config.SchemaVersion,
		RulePacks: schemaRulePackConfig{
			Enabled: append([]string{}, config.RulePacks.Enabled...),
		},
		Repository: repositoryToSchema(config.Repository),
		StyleGuide: schemaStyleGuideConfig{
			Path:                config.StyleGuide.Path,
			RequirementIDFormat: config.StyleGuide.RequirementIDFormat,
		},
		Formatting:   formattingToSchema(config.Formatting),
		Imports:      schemaImportsConfig{LocalPrefix: config.Imports.LocalPrefix},
		Paths:        cloneStringMap(config.Paths.Classes),
		FileSets:     fileSetsToSchema(config.FileSets),
		Language:     languageToSchema(config.Language),
		Tools:        toolsToSchema(config.Tools),
		Naming:       namingToSchema(config.Naming),
		ControlPlane: controlPlaneToSchema(config.ControlPlane),
		Architecture: architectureToSchema(config.Architecture),
		Rules:        rulesToSchema(config.Rules),
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

func cloneStringMap[M ~map[string][]string](source M) (target M) {
	if source == nil {
		return nil
	}

	target = make(M, len(source))
	for key, values := range source {
		target[key] = append([]string{}, values...)
	}

	return target
}
