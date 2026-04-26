package profile

import "ciphera/tools/internal/policy"

type schemaImportsConfig struct {
	LocalPrefix string `toml:"local_prefix"`
}

type schemaLanguageConfig struct {
	Backends []schemaLanguageBackend `toml:"backends"`
}

type schemaLanguageBackend struct {
	Name        string   `toml:"name"`
	Language    string   `toml:"language"`
	Scope       string   `toml:"scope"`
	Workdir     string   `toml:"workdir"`
	FormatPaths []string `toml:"format_paths"`
	StylePaths  []string `toml:"style_paths"`
}

type schemaToolPin struct {
	ID               string `toml:"id"`
	Version          string `toml:"version"`
	TimeoutSeconds   int    `toml:"timeout_seconds"`
	OutputLimitBytes int64  `toml:"output_limit_bytes"`
}

type schemaNamingConfig struct {
	TypeSuffixForbidden []string                        `toml:"go_type_suffix_forbidden"`
	TypeSuffixPreferred string                          `toml:"go_type_suffix_preferred"`
	IdentifierForbidden []string                        `toml:"go_identifier_suffix_forbidden"`
	IdentifierPreferred string                          `toml:"go_identifier_suffix_preferred"`
	GoParameters        schemaGoParameterConfig         `toml:"go_parameters"`
	GoDomainIdentifiers policy.GoDomainIdentifierConfig `toml:"go_domain_identifiers"`
	ShellForbidden      []string                        `toml:"shell_forbidden_assignments"`
	ShellPreferred      string                          `toml:"shell_preferred_assignment"`
}

type schemaGoParameterConfig struct {
	SecretNames           []string                      `toml:"secret_names"`
	ConstructorCategories []schemaGoConstructorCategory `toml:"constructor_categories"`
}

type schemaGoConstructorCategory struct {
	Name                string   `toml:"name"`
	TypeMarkers         []string `toml:"type_markers"`
	ExcludedTypeMarkers []string `toml:"excluded_type_markers"`
	ParameterNames      []string `toml:"parameter_names"`
	UsesSecretNames     bool     `toml:"uses_secret_names"`
}

type schemaControlPlane struct {
	QualityFile       string                       `toml:"quality_file"`
	VariableContracts []schemaMakeVariableContract `toml:"variable_contracts"`
	TargetContracts   []schemaMakeTargetContract   `toml:"target_contracts"`
}

type schemaMakeVariableContract struct {
	Name  string `toml:"name"`
	Value string `toml:"value"`
}

type schemaMakeTargetContract struct {
	Name       string `toml:"name"`
	RecipeLine string `toml:"recipe_line"`
}

type schemaArchitecture struct {
	Layers []schemaArchitectureLayer `toml:"layers"`
}

type schemaArchitectureLayer struct {
	Name         string   `toml:"name"`
	PackageRoots []string `toml:"package_roots"`
	MayImport    []string `toml:"may_import"`
}
