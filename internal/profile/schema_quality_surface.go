package profile

type schemaQualitySurface struct {
	Driver string           `toml:"driver"`
	Make   schemaMakeConfig `toml:"make"`
}

type schemaMakeConfig struct {
	Path              string                   `toml:"path"`
	RequiredVariables []schemaMakefileVariable `toml:"required_variables"`
	RequiredTargets   []schemaMakefileTarget   `toml:"required_targets"`
}

type schemaMakefileVariable struct {
	Name  string `toml:"name"`
	Value string `toml:"value"`
}

type schemaMakefileTarget struct {
	Name       string `toml:"name"`
	RecipeLine string `toml:"recipe_line"`
}
