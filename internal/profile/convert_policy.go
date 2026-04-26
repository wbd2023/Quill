package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ------------------------------------------ Language ------------------------------------------ */

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

/* -------------------------------------------- Tools ------------------------------------------- */

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

/* ------------------------------------------- Naming ------------------------------------------- */

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

/* ---------------------------------------- Control Plane --------------------------------------- */

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

/* ---------------------------------------- Architecture ---------------------------------------- */

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
