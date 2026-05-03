package profile

import "ciphera/tools/internal/policy"

type schemaGoConfig struct {
	LocalImportPrefixes    []string                              `toml:"local_import_prefixes"`
	Parameters             schemaGoParameterConfig               `toml:"parameters"`
	IdentifierConstructors policy.GoDomainIdentifierConstructors `toml:"domain_identifiers"`
	Architecture           schemaGoArchitecture                  `toml:"architecture"`
}

type schemaGoParameterConfig struct {
	SecretNames      []string               `toml:"secret_names"`
	ConstructorOrder []schemaParameterGroup `toml:"constructor_order"`
}

type schemaParameterGroup struct {
	Name                string   `toml:"name"`
	TypeMarkers         []string `toml:"type_markers"`
	ExcludedTypeMarkers []string `toml:"excluded_type_markers"`
	ParameterNames      []string `toml:"parameter_names"`
	MatchesSecretNames  bool     `toml:"matches_secret_names"`
}

type schemaGoArchitecture struct {
	Layers []schemaGoArchitectureLayer `toml:"layers"`
}

type schemaGoArchitectureLayer struct {
	Name          string   `toml:"name"`
	PackageRoots  []string `toml:"package_roots"`
	AllowedLayers []string `toml:"may_import"`
}
