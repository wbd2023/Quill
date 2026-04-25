package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/contract"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const SchemaVersion1 = 1

const RequirementIDFormatSectionSlug = "section_slug"

/* -------------------------------------------- Types ------------------------------------------- */

type Profile struct {
	SchemaVersion int                `toml:"profile_version"`
	RulePacks     RulePackConfig     `toml:"rule_packs"`
	Repository    RepositoryConfig   `toml:"repository"`
	StyleGuide    StyleGuideConfig   `toml:"styleguide"`
	Imports       ImportsConfig      `toml:"imports"`
	Paths         PathClassSet       `toml:"paths"`
	FileSets      []FileSetConfig    `toml:"file_sets"`
	Language      LanguageConfig     `toml:"language"`
	Naming        NamingConfig       `toml:"naming"`
	ControlPlane  ControlPlaneConfig `toml:"control_plane"`
	Architecture  ArchitectureConfig `toml:"architecture"`
	Rules         []RuleBinding      `toml:"rules"`
}

type RulePackConfig struct {
	Enabled []string `toml:"enabled"`
}

type RepositoryConfig struct {
	RootMarkers         []string `toml:"root_markers"`
	AppScanRoots        []string `toml:"app_scan_roots"`
	ToolsScanRoots      []string `toml:"tools_scan_roots"`
	GlobalExclusions    []string `toml:"global_exclusions"`
	GeneratedMarker     string   `toml:"generated_marker"`
	GeneratedProbeLimit int      `toml:"generated_probe_limit"`
}

type StyleGuideConfig struct {
	Path                string `toml:"path"`
	RequirementIDFormat string `toml:"requirement_id_format"`
}

type ImportsConfig struct {
	LocalPrefix string `toml:"local_prefix"`
}

type PathClassSet struct {
	Classes map[string][]string
}

type FileSetConfig struct {
	Name                 string   `toml:"name"`
	Extensions           []string `toml:"extensions"`
	AppFiles             []string `toml:"app_files"`
	AppPrefixes          []string `toml:"app_prefixes"`
	ToolsFiles           []string `toml:"tools_files"`
	ToolsPrefixes        []string `toml:"tools_prefixes"`
	ExcludedExtensions   []string `toml:"excluded_extensions"`
	ExcludedNames        []string `toml:"excluded_names"`
	ExcludedNamePrefixes []string `toml:"excluded_name_prefixes"`
	SkipBinary           bool     `toml:"skip_binary"`
}

type LanguageConfig struct {
	Backends []LanguageBackendConfig `toml:"backends"`
}

type LanguageBackendConfig struct {
	Name        string   `toml:"name"`
	Language    string   `toml:"language"`
	Workdir     string   `toml:"workdir"`
	FormatPaths []string `toml:"format_paths"`
	StylePaths  []string `toml:"style_paths"`
}

type NamingConfig struct {
	GoTypeSuffixForbidden       []string                 `toml:"go_type_suffix_forbidden"`
	GoTypeSuffixPreferred       string                   `toml:"go_type_suffix_preferred"`
	GoIdentifierSuffixForbidden []string                 `toml:"go_identifier_suffix_forbidden"`
	GoIdentifierSuffixPreferred string                   `toml:"go_identifier_suffix_preferred"`
	GoParameters                GoParameterConfig        `toml:"go_parameters"`
	GoDomainIdentifiers         GoDomainIdentifierConfig `toml:"go_domain_identifiers"`
	ShellForbiddenAssignments   []string                 `toml:"shell_forbidden_assignments"`
	ShellPreferredAssignment    string                   `toml:"shell_preferred_assignment"`
}

type GoParameterConfig struct {
	SecretNames           []string                `toml:"secret_names"`
	ConstructorCategories []GoConstructorCategory `toml:"constructor_categories"`
}

type GoConstructorCategory struct {
	Name                string   `toml:"name"`
	TypeMarkers         []string `toml:"type_markers"`
	ExcludedTypeMarkers []string `toml:"excluded_type_markers"`
	ParameterNames      []string `toml:"parameter_names"`
	UsesSecretNames     bool     `toml:"uses_secret_names"`
}

type GoDomainIdentifierConfig map[string][]string

type ControlPlaneConfig struct {
	QualityFile       string                          `toml:"quality_file"`
	VariableContracts []contract.MakeVariableContract `toml:"variable_contracts"`
	TargetContracts   []contract.MakeTargetContract   `toml:"target_contracts"`
}

type ArchitectureConfig struct {
	Layers []ArchitectureLayer `toml:"layers"`
}

type ArchitectureLayer struct {
	Name         string   `toml:"name"`
	PackageRoots []string `toml:"package_roots"`
	MayImport    []string `toml:"may_import"`
}

type RuleBinding struct {
	RuleID         string         `toml:"rule_id"`
	Level          contract.Level `toml:"level"`
	Scope          contract.Scope `toml:"scope"`
	RequirementIDs []string       `toml:"requirement_ids"`
	ConfigRef      string         `toml:"config_ref"`
}

type EffectiveRule = contract.Rule

type EffectiveConfig struct {
	Tools []contract.Tool
	Rules []EffectiveRule
}

func (effective EffectiveConfig) ToolByID(id string) (tool contract.Tool, found bool) {
	for _, current := range effective.Tools {
		if current.ID == id {
			return current, true
		}
	}

	return contract.Tool{}, false
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func (repository RepositoryConfig) ScanRoots(
	repoRoot string,
	scope contract.Scope,
) (roots []string) {
	switch scope {
	case contract.ScopeApp:
		return joinPaths(repoRoot, repository.AppScanRoots)

	case contract.ScopeTools:
		return joinPaths(repoRoot, repository.ToolsScanRoots)

	case contract.ScopeAll:
		return []string{repoRoot}

	default:
		return nil
	}
}

func (repository RepositoryConfig) ValidateRoot(repoRoot string) (err error) {
	for _, marker := range repository.RootMarkers {
		if marker == "" {
			continue
		}

		if _, statErr := os.Stat(filepath.Join(repoRoot, marker)); statErr == nil {
			continue
		} else {
			return statErr
		}
	}

	return nil
}

func (paths PathClassSet) Patterns(className string) (patterns []string) {
	if paths.Classes == nil {
		return nil
	}

	return append([]string{}, paths.Classes[className]...)
}

func (paths *PathClassSet) UnmarshalTOML(value any) (err error) {
	rawClasses, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("paths must be a table")
	}

	paths.Classes = make(map[string][]string, len(rawClasses))
	for className, rawValues := range rawClasses {
		values, err := stringList(rawValues)
		if err != nil {
			return fmt.Errorf("paths.%s: %w", className, err)
		}

		paths.Classes[className] = values
	}

	return nil
}

func stringList(value any) (values []string, err error) {
	switch typed := value.(type) {
	case []string:
		return append([]string{}, typed...), nil

	case []any:
		values = make([]string, 0, len(typed))
		for _, rawItem := range typed {
			item, ok := rawItem.(string)
			if !ok {
				return nil, fmt.Errorf("must contain only strings")
			}

			values = append(values, item)
		}

		return values, nil

	default:
		return nil, fmt.Errorf("must be a string array")
	}
}

func (policy Profile) FileSet(name string) (fileSet FileSetConfig, found bool) {
	for _, current := range policy.FileSets {
		if current.Name == name {
			return current, true
		}
	}

	return FileSetConfig{}, false
}

func (policy Profile) LanguageBackend(
	name string,
) (backend LanguageBackendConfig, found bool) {
	for _, current := range policy.Language.Backends {
		if current.Name == name {
			return current, true
		}
	}

	return LanguageBackendConfig{}, false
}

func joinPaths(repoRoot string, values []string) (paths []string) {
	paths = make([]string, 0, len(values))

	for _, value := range values {
		if value == "." {
			paths = append(paths, repoRoot)
			continue
		}

		paths = append(paths, filepath.Join(repoRoot, value))
	}

	return paths
}
