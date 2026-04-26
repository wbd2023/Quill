package profile

import "ciphera/tools/internal/policy"

/* ----------------------------------------- Root Schema ---------------------------------------- */

type schemaConfig struct {
	SchemaVersion int                    `toml:"profile_version"`
	RulePacks     schemaRulePackConfig   `toml:"rule_packs"`
	Repository    schemaRepositoryConfig `toml:"repository"`
	StyleGuide    schemaStyleGuideConfig `toml:"styleguide"`
	Formatting    schemaFormattingConfig `toml:"formatting"`
	Imports       schemaImportsConfig    `toml:"imports"`
	Paths         map[string][]string    `toml:"paths"`
	FileSets      []schemaFileSetConfig  `toml:"file_sets"`
	Language      schemaLanguageConfig   `toml:"language"`
	Tools         []schemaToolPin        `toml:"tools"`
	Naming        schemaNamingConfig     `toml:"naming"`
	ControlPlane  schemaControlPlane     `toml:"control_plane"`
	Architecture  schemaArchitecture     `toml:"architecture"`
	Rules         []schemaRuleBinding    `toml:"rules"`
}

type schemaRulePackConfig struct {
	Enabled []string `toml:"enabled"`
}

type schemaRepositoryConfig struct {
	RootMarkers         []string            `toml:"root_markers"`
	DefaultScope        string              `toml:"default_scope"`
	Scopes              map[string][]string `toml:"scopes"`
	GlobalExclusions    []string            `toml:"global_exclusions"`
	GeneratedMarker     string              `toml:"generated_marker"`
	GeneratedProbeLimit int                 `toml:"generated_probe_limit"`
}

type schemaStyleGuideConfig struct {
	Path                string `toml:"path"`
	RequirementIDFormat string `toml:"requirement_id_format"`
}

type schemaFormattingConfig struct {
	SectionHeaders schemaSectionHeaderConfig `toml:"section_headers"`
}

type schemaSectionHeaderConfig struct {
	RequiredMinLines  int      `toml:"required_min_lines"`
	ShortFileMaxLines int      `toml:"short_file_max_lines"`
	OveruseCount      int      `toml:"overuse_header_count"`
	GenericNames      []string `toml:"generic_names"`
	StructuralNames   []string `toml:"structural_names"`
}

type schemaImportsConfig struct {
	LocalPrefix string `toml:"local_prefix"`
}

/* ----------------------------------------- File Schema ---------------------------------------- */

type schemaFileSetConfig struct {
	Name                 string              `toml:"name"`
	Extensions           []string            `toml:"extensions"`
	Files                map[string][]string `toml:"files"`
	Prefixes             map[string][]string `toml:"prefixes"`
	ExcludedExtensions   []string            `toml:"excluded_extensions"`
	ExcludedNames        []string            `toml:"excluded_names"`
	ExcludedNamePrefixes []string            `toml:"excluded_name_prefixes"`
	SkipBinary           bool                `toml:"skip_binary"`
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

/* ---------------------------------------- Naming Schema --------------------------------------- */

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

/* --------------------------------------- Control Schema --------------------------------------- */

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

/* ------------------------------------- Architecture Schema ------------------------------------ */

type schemaArchitecture struct {
	Layers []schemaArchitectureLayer `toml:"layers"`
}

type schemaArchitectureLayer struct {
	Name         string   `toml:"name"`
	PackageRoots []string `toml:"package_roots"`
	MayImport    []string `toml:"may_import"`
}
