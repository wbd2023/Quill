package profile

import "ciphera/tools/internal/policy"

func goFromSchema(schema schemaGoConfig) (config policy.GoConfig) {
	return policy.GoConfig{
		LocalImportPrefixes:          append([]string{}, schema.LocalImportPrefixes...),
		Parameters:                   goParametersFromSchema(schema.Parameters),
		DomainIdentifierConstructors: cloneStringSlices(schema.IdentifierConstructors),
		Architecture:                 goArchitectureFromSchema(schema.Architecture),
	}
}

func goToSchema(goConfig policy.GoConfig) (schema schemaGoConfig) {
	return schemaGoConfig{
		LocalImportPrefixes:    append([]string{}, goConfig.LocalImportPrefixes...),
		Parameters:             goParametersToSchema(goConfig.Parameters),
		IdentifierConstructors: cloneStringSlices(goConfig.DomainIdentifierConstructors),
		Architecture:           goArchitectureToSchema(goConfig.Architecture),
	}
}

func goParametersFromSchema(schema schemaGoParameterConfig) (parameters policy.GoParameterConfig) {
	parameters.SecretNames = append([]string{}, schema.SecretNames...)
	parameters.ConstructorOrder = make(
		[]policy.GoParameterGroup,
		0,
		len(schema.ConstructorOrder),
	)
	for _, group := range schema.ConstructorOrder {
		parameters.ConstructorOrder = append(
			parameters.ConstructorOrder,
			policy.GoParameterGroup{
				Name:                group.Name,
				TypeMarkers:         append([]string{}, group.TypeMarkers...),
				ExcludedTypeMarkers: append([]string{}, group.ExcludedTypeMarkers...),
				ParameterNames:      append([]string{}, group.ParameterNames...),
				MatchesSecretNames:  group.MatchesSecretNames,
			},
		)
	}

	return parameters
}

func goParametersToSchema(parameters policy.GoParameterConfig) (schema schemaGoParameterConfig) {
	schema.SecretNames = append([]string{}, parameters.SecretNames...)
	schema.ConstructorOrder = make(
		[]schemaParameterGroup,
		0,
		len(parameters.ConstructorOrder),
	)
	for _, group := range parameters.ConstructorOrder {
		schema.ConstructorOrder = append(
			schema.ConstructorOrder,
			schemaParameterGroup{
				Name:                group.Name,
				TypeMarkers:         append([]string{}, group.TypeMarkers...),
				ExcludedTypeMarkers: append([]string{}, group.ExcludedTypeMarkers...),
				ParameterNames:      append([]string{}, group.ParameterNames...),
				MatchesSecretNames:  group.MatchesSecretNames,
			},
		)
	}

	return schema
}

func goArchitectureFromSchema(
	schema schemaGoArchitecture,
) (architecture policy.GoArchitectureConfig) {
	architecture.Layers = make([]policy.GoArchitectureLayer, 0, len(schema.Layers))
	for _, layer := range schema.Layers {
		architecture.Layers = append(architecture.Layers, policy.GoArchitectureLayer{
			Name:          layer.Name,
			PackageRoots:  append([]string{}, layer.PackageRoots...),
			AllowedLayers: append([]string{}, layer.AllowedLayers...),
		})
	}

	return architecture
}

func goArchitectureToSchema(
	architecture policy.GoArchitectureConfig,
) (schema schemaGoArchitecture) {
	schema.Layers = make([]schemaGoArchitectureLayer, 0, len(architecture.Layers))
	for _, layer := range architecture.Layers {
		schema.Layers = append(schema.Layers, schemaGoArchitectureLayer{
			Name:          layer.Name,
			PackageRoots:  append([]string{}, layer.PackageRoots...),
			AllowedLayers: append([]string{}, layer.AllowedLayers...),
		})
	}

	return schema
}
