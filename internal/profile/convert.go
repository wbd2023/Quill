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

func repositoryFromSchema(schema schemaRepositoryConfig) (repository policy.RepositoryConfig) {
	return policy.RepositoryConfig{
		RootMarkers:         append([]string{}, schema.RootMarkers...),
		DefaultScope:        contract.Scope(schema.DefaultScope),
		Scopes:              scopeMapFromSchema(schema.Scopes),
		GlobalExclusions:    append([]string{}, schema.GlobalExclusions...),
		GeneratedMarker:     schema.GeneratedMarker,
		GeneratedProbeLimit: schema.GeneratedProbeLimit,
	}
}

func repositoryToSchema(repository policy.RepositoryConfig) (schema schemaRepositoryConfig) {
	return schemaRepositoryConfig{
		RootMarkers:         append([]string{}, repository.RootMarkers...),
		DefaultScope:        string(repository.DefaultScope),
		Scopes:              scopeMapToSchema(repository.Scopes),
		GlobalExclusions:    append([]string{}, repository.GlobalExclusions...),
		GeneratedMarker:     repository.GeneratedMarker,
		GeneratedProbeLimit: repository.GeneratedProbeLimit,
	}
}

func styleGuideFromSchema(schema schemaStyleGuideConfig) (config policy.StyleGuideConfig) {
	return policy.StyleGuideConfig{
		Path:                schema.Path,
		RequirementIDFormat: schema.RequirementIDFormat,
	}
}

func formattingFromSchema(schema schemaFormattingConfig) (config policy.FormattingConfig) {
	return policy.FormattingConfig{
		SectionHeaders: policy.SectionHeaderConfig{
			RequiredMinLines:  schema.SectionHeaders.RequiredMinLines,
			ShortFileMaxLines: schema.SectionHeaders.ShortFileMaxLines,
			OveruseCount:      schema.SectionHeaders.OveruseCount,
			GenericNames:      append([]string{}, schema.SectionHeaders.GenericNames...),
			StructuralNames:   append([]string{}, schema.SectionHeaders.StructuralNames...),
		},
	}
}

func formattingToSchema(config policy.FormattingConfig) (schema schemaFormattingConfig) {
	return schemaFormattingConfig{
		SectionHeaders: schemaSectionHeaderConfig{
			RequiredMinLines:  config.SectionHeaders.RequiredMinLines,
			ShortFileMaxLines: config.SectionHeaders.ShortFileMaxLines,
			OveruseCount:      config.SectionHeaders.OveruseCount,
			GenericNames:      append([]string{}, config.SectionHeaders.GenericNames...),
			StructuralNames:   append([]string{}, config.SectionHeaders.StructuralNames...),
		},
	}
}

func fileSetsFromSchema(schemas []schemaFileSetConfig) (fileSets []policy.FileSetConfig) {
	fileSets = make([]policy.FileSetConfig, 0, len(schemas))
	for _, fileSet := range schemas {
		fileSets = append(fileSets, policy.FileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			Files:                scopeMapFromSchema(fileSet.Files),
			Prefixes:             scopeMapFromSchema(fileSet.Prefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return fileSets
}

func fileSetsToSchema(fileSets []policy.FileSetConfig) (schemas []schemaFileSetConfig) {
	schemas = make([]schemaFileSetConfig, 0, len(fileSets))
	for _, fileSet := range fileSets {
		schemas = append(schemas, schemaFileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			Files:                scopeMapToSchema(fileSet.Files),
			Prefixes:             scopeMapToSchema(fileSet.Prefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return schemas
}

func languageFromSchema(schema schemaLanguageConfig) (language policy.LanguageConfig) {
	language.Backends = make([]policy.LanguageBackendConfig, 0, len(schema.Backends))
	for _, backend := range schema.Backends {
		language.Backends = append(language.Backends, policy.LanguageBackendConfig{
			Name:        backend.Name,
			Language:    backend.Language,
			Scope:       contract.Scope(backend.Scope),
			Workdir:     backend.Workdir,
			FormatPaths: append([]string{}, backend.FormatPaths...),
			StylePaths:  append([]string{}, backend.StylePaths...),
		})
	}

	return language
}

func languageToSchema(language policy.LanguageConfig) (schema schemaLanguageConfig) {
	schema.Backends = make([]schemaLanguageBackend, 0, len(language.Backends))
	for _, backend := range language.Backends {
		schema.Backends = append(schema.Backends, schemaLanguageBackend{
			Name:        backend.Name,
			Language:    backend.Language,
			Scope:       string(backend.Scope),
			Workdir:     backend.Workdir,
			FormatPaths: append([]string{}, backend.FormatPaths...),
			StylePaths:  append([]string{}, backend.StylePaths...),
		})
	}

	return schema
}

func toolsFromSchema(schemas []schemaToolPin) (tools []policy.ToolPin) {
	tools = make([]policy.ToolPin, 0, len(schemas))
	for _, tool := range schemas {
		tools = append(tools, policy.ToolPin{
			ID:               tool.ID,
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		})
	}

	return tools
}

func toolsToSchema(tools []policy.ToolPin) (schemas []schemaToolPin) {
	schemas = make([]schemaToolPin, 0, len(tools))
	for _, tool := range tools {
		schemas = append(schemas, schemaToolPin{
			ID:               tool.ID,
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		})
	}

	return schemas
}

func namingFromSchema(schema schemaNamingConfig) (naming policy.NamingConfig) {
	return policy.NamingConfig{
		GoTypeSuffixForbidden:       append([]string{}, schema.TypeSuffixForbidden...),
		GoTypeSuffixPreferred:       schema.TypeSuffixPreferred,
		GoIdentifierSuffixForbidden: append([]string{}, schema.IdentifierForbidden...),
		GoIdentifierSuffixPreferred: schema.IdentifierPreferred,
		GoParameters:                goParametersFromSchema(schema.GoParameters),
		GoDomainIdentifiers:         cloneStringMap(schema.GoDomainIdentifiers),
		ShellForbiddenAssignments:   append([]string{}, schema.ShellForbidden...),
		ShellPreferredAssignment:    schema.ShellPreferred,
	}
}

func namingToSchema(naming policy.NamingConfig) (schema schemaNamingConfig) {
	return schemaNamingConfig{
		TypeSuffixForbidden: append([]string{}, naming.GoTypeSuffixForbidden...),
		TypeSuffixPreferred: naming.GoTypeSuffixPreferred,
		IdentifierForbidden: append([]string{}, naming.GoIdentifierSuffixForbidden...),
		IdentifierPreferred: naming.GoIdentifierSuffixPreferred,
		GoParameters:        goParametersToSchema(naming.GoParameters),
		GoDomainIdentifiers: cloneStringMap(naming.GoDomainIdentifiers),
		ShellForbidden:      append([]string{}, naming.ShellForbiddenAssignments...),
		ShellPreferred:      naming.ShellPreferredAssignment,
	}
}

func goParametersFromSchema(schema schemaGoParameterConfig) (parameters policy.GoParameterConfig) {
	parameters.SecretNames = append([]string{}, schema.SecretNames...)
	parameters.ConstructorCategories = make(
		[]policy.GoConstructorCategory,
		0,
		len(schema.ConstructorCategories),
	)
	for _, category := range schema.ConstructorCategories {
		parameters.ConstructorCategories = append(
			parameters.ConstructorCategories,
			policy.GoConstructorCategory{
				Name:                category.Name,
				TypeMarkers:         append([]string{}, category.TypeMarkers...),
				ExcludedTypeMarkers: append([]string{}, category.ExcludedTypeMarkers...),
				ParameterNames:      append([]string{}, category.ParameterNames...),
				UsesSecretNames:     category.UsesSecretNames,
			},
		)
	}

	return parameters
}

func goParametersToSchema(parameters policy.GoParameterConfig) (schema schemaGoParameterConfig) {
	schema.SecretNames = append([]string{}, parameters.SecretNames...)
	schema.ConstructorCategories = make(
		[]schemaGoConstructorCategory,
		0,
		len(parameters.ConstructorCategories),
	)
	for _, category := range parameters.ConstructorCategories {
		schema.ConstructorCategories = append(
			schema.ConstructorCategories,
			schemaGoConstructorCategory{
				Name:                category.Name,
				TypeMarkers:         append([]string{}, category.TypeMarkers...),
				ExcludedTypeMarkers: append([]string{}, category.ExcludedTypeMarkers...),
				ParameterNames:      append([]string{}, category.ParameterNames...),
				UsesSecretNames:     category.UsesSecretNames,
			},
		)
	}

	return schema
}

func controlPlaneFromSchema(schema schemaControlPlane) (control policy.ControlPlaneConfig) {
	control.QualityFile = schema.QualityFile
	control.VariableContracts = make(
		[]policy.MakeVariableContract,
		0,
		len(schema.VariableContracts),
	)
	for _, contract := range schema.VariableContracts {
		control.VariableContracts = append(control.VariableContracts, policy.MakeVariableContract{
			Name:  contract.Name,
			Value: contract.Value,
		})
	}

	control.TargetContracts = make([]policy.MakeTargetContract, 0, len(schema.TargetContracts))
	for _, contract := range schema.TargetContracts {
		control.TargetContracts = append(control.TargetContracts, policy.MakeTargetContract{
			Name:       contract.Name,
			RecipeLine: contract.RecipeLine,
		})
	}

	return control
}

func controlPlaneToSchema(control policy.ControlPlaneConfig) (schema schemaControlPlane) {
	schema.QualityFile = control.QualityFile
	schema.VariableContracts = make(
		[]schemaMakeVariableContract,
		0,
		len(control.VariableContracts),
	)
	for _, contract := range control.VariableContracts {
		schema.VariableContracts = append(schema.VariableContracts, schemaMakeVariableContract{
			Name:  contract.Name,
			Value: contract.Value,
		})
	}

	schema.TargetContracts = make([]schemaMakeTargetContract, 0, len(control.TargetContracts))
	for _, contract := range control.TargetContracts {
		schema.TargetContracts = append(schema.TargetContracts, schemaMakeTargetContract{
			Name:       contract.Name,
			RecipeLine: contract.RecipeLine,
		})
	}

	return schema
}

func architectureFromSchema(schema schemaArchitecture) (architecture policy.ArchitectureConfig) {
	architecture.Layers = make([]policy.ArchitectureLayer, 0, len(schema.Layers))
	for _, layer := range schema.Layers {
		architecture.Layers = append(architecture.Layers, policy.ArchitectureLayer{
			Name:         layer.Name,
			PackageRoots: append([]string{}, layer.PackageRoots...),
			MayImport:    append([]string{}, layer.MayImport...),
		})
	}

	return architecture
}

func architectureToSchema(architecture policy.ArchitectureConfig) (schema schemaArchitecture) {
	schema.Layers = make([]schemaArchitectureLayer, 0, len(architecture.Layers))
	for _, layer := range architecture.Layers {
		schema.Layers = append(schema.Layers, schemaArchitectureLayer{
			Name:         layer.Name,
			PackageRoots: append([]string{}, layer.PackageRoots...),
			MayImport:    append([]string{}, layer.MayImport...),
		})
	}

	return schema
}

func rulesFromSchema(schemas []schemaRuleBinding) (rules []policy.RuleBinding) {
	rules = make([]policy.RuleBinding, 0, len(schemas))
	for _, rule := range schemas {
		rules = append(rules, policy.RuleBinding{
			RuleID:         rule.RuleID,
			Level:          rule.Level,
			Scope:          contract.Scope(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
			ConfigRef:      rule.ConfigRef,
			Backends:       append([]string{}, rule.Backends...),
			PathClasses:    append([]string{}, rule.PathClasses...),
		})
	}

	return rules
}

func rulesToSchema(rules []policy.RuleBinding) (schemas []schemaRuleBinding) {
	schemas = make([]schemaRuleBinding, 0, len(rules))
	for _, rule := range rules {
		schemas = append(schemas, schemaRuleBinding{
			RuleID:         rule.RuleID,
			Level:          rule.Level,
			Scope:          string(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
			ConfigRef:      rule.ConfigRef,
			Backends:       append([]string{}, rule.Backends...),
			PathClasses:    append([]string{}, rule.PathClasses...),
		})
	}

	return schemas
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
