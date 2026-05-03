package profile

import "ciphera/tools/internal/policy"

func qualitySurfaceFromSchema(schema schemaQualitySurface) (surface policy.QualitySurfaceConfig) {
	surface.Driver = policy.QualitySurfaceDriver(schema.Driver)
	surface.Make.Path = schema.Make.Path
	surface.Make.RequiredVariables = make(
		[]policy.MakefileVariable,
		0,
		len(schema.Make.RequiredVariables),
	)
	for _, variable := range schema.Make.RequiredVariables {
		surface.Make.RequiredVariables = append(
			surface.Make.RequiredVariables,
			policy.MakefileVariable{
				Name:  variable.Name,
				Value: variable.Value,
			},
		)
	}

	surface.Make.RequiredTargets = make(
		[]policy.MakefileTarget,
		0,
		len(schema.Make.RequiredTargets),
	)
	for _, target := range schema.Make.RequiredTargets {
		surface.Make.RequiredTargets = append(surface.Make.RequiredTargets, policy.MakefileTarget{
			Name:       target.Name,
			RecipeLine: target.RecipeLine,
		})
	}

	return surface
}

func qualitySurfaceToSchema(surface policy.QualitySurfaceConfig) (schema schemaQualitySurface) {
	schema.Driver = string(surface.Driver)
	schema.Make.Path = surface.Make.Path
	schema.Make.RequiredVariables = make(
		[]schemaMakefileVariable,
		0,
		len(surface.Make.RequiredVariables),
	)
	for _, variable := range surface.Make.RequiredVariables {
		schema.Make.RequiredVariables = append(
			schema.Make.RequiredVariables,
			schemaMakefileVariable{
				Name:  variable.Name,
				Value: variable.Value,
			},
		)
	}

	schema.Make.RequiredTargets = make(
		[]schemaMakefileTarget,
		0,
		len(surface.Make.RequiredTargets),
	)
	for _, target := range surface.Make.RequiredTargets {
		schema.Make.RequiredTargets = append(schema.Make.RequiredTargets, schemaMakefileTarget{
			Name:       target.Name,
			RecipeLine: target.RecipeLine,
		})
	}

	return schema
}
