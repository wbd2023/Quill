package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func Validate(config policy.Config) (err error) {
	if config.SchemaVersion != policy.SchemaVersion {
		return fmt.Errorf("unsupported style profile version %d", config.SchemaVersion)
	}

	if len(config.RulePacks.Enabled) == 0 {
		return fmt.Errorf("rule_packs.enabled must not be empty")
	}

	if len(config.Repository.RootMarkers) == 0 {
		return fmt.Errorf("repository.root_markers must not be empty")
	}

	if len(config.Repository.Scopes) == 0 {
		return fmt.Errorf("repository.scopes must not be empty")
	}

	if config.Repository.DefaultScope == "" {
		return fmt.Errorf("repository.default_scope must not be empty")
	}

	if config.Repository.GeneratedMarker == "" {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	if config.Repository.GeneratedProbeLimit <= 0 {
		return fmt.Errorf("repository.generated_probe_limit must be positive")
	}

	if config.StyleGuide.Path == "" {
		return fmt.Errorf("styleguide.path must not be empty")
	}

	if config.StyleGuide.RequirementIDFormat == "" {
		return fmt.Errorf("styleguide.requirement_id_format must not be empty")
	}

	if config.StyleGuide.RequirementIDFormat != policy.RequirementIDFormatSectionSlug {
		return fmt.Errorf(
			"unsupported styleguide.requirement_id_format %q",
			config.StyleGuide.RequirementIDFormat,
		)
	}

	if err = validateFormatting(config.Formatting); err != nil {
		return err
	}

	if config.Imports.LocalPrefix == "" {
		return fmt.Errorf("imports.local_prefix must not be empty")
	}

	for scope, roots := range config.Repository.Scopes {
		if scope == "" {
			return fmt.Errorf("repository.scopes contains an empty scope")
		}

		if len(roots) == 0 {
			return fmt.Errorf("repository.scopes.%s must not be empty", scope)
		}
	}

	if !config.Repository.ScopeExists(config.Repository.DefaultScope) {
		return fmt.Errorf(
			"repository.default_scope references unknown scope %q",
			config.Repository.DefaultScope,
		)
	}

	if err = validateGoDomainIdentifiers(config.Naming.GoDomainIdentifiers); err != nil {
		return err
	}

	seenFileSets := make(map[string]bool, len(config.FileSets))
	for _, fileSet := range config.FileSets {
		if fileSet.Name == "" {
			return fmt.Errorf("file set name must not be empty")
		}

		if seenFileSets[fileSet.Name] {
			return fmt.Errorf("duplicate file set %q", fileSet.Name)
		}

		seenFileSets[fileSet.Name] = true
		for scope := range fileSet.Files {
			if !config.Repository.ScopeExists(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}

		for scope := range fileSet.Prefixes {
			if !config.Repository.ScopeExists(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}
	}

	seenLanguageBackends := make(map[string]bool, len(config.Language.Backends))
	for _, backend := range config.Language.Backends {
		if backend.Name == "" {
			return fmt.Errorf("language backend name must not be empty")
		}

		if seenLanguageBackends[backend.Name] {
			return fmt.Errorf("duplicate language backend %q", backend.Name)
		}

		seenLanguageBackends[backend.Name] = true

		if backend.Language == "" {
			return fmt.Errorf("language backend %q must define language", backend.Name)
		}

		if !config.Repository.ScopeExists(backend.Scope) {
			return fmt.Errorf(
				"language backend %q references unknown scope %q",
				backend.Name,
				backend.Scope,
			)
		}
	}

	if config.ControlPlane.QualityFile == "" {
		return fmt.Errorf("control_plane.quality_file must not be empty")
	}

	if len(config.Architecture.Layers) == 0 {
		return fmt.Errorf("architecture.layers must not be empty")
	}

	if len(config.Rules) == 0 {
		return fmt.Errorf("rules must not be empty")
	}

	layerNames := make(map[string]bool, len(config.Architecture.Layers))
	for _, layer := range config.Architecture.Layers {
		if layer.Name == "" {
			return fmt.Errorf("architecture layer name must not be empty")
		}

		if layerNames[layer.Name] {
			return fmt.Errorf("duplicate architecture layer %q", layer.Name)
		}

		layerNames[layer.Name] = true

		if len(layer.PackageRoots) == 0 {
			return fmt.Errorf("architecture layer %q must define package_roots", layer.Name)
		}
	}

	for _, layer := range config.Architecture.Layers {
		for _, allowed := range layer.MayImport {
			if layerNames[allowed] {
				continue
			}

			return fmt.Errorf(
				"architecture layer %q references unknown layer %q",
				layer.Name,
				allowed,
			)
		}
	}

	seenToolPins := make(map[string]bool, len(config.Tools))
	for _, tool := range config.Tools {
		if tool.ID == "" {
			return fmt.Errorf("tool pin has an empty id")
		}

		if seenToolPins[tool.ID] {
			return fmt.Errorf("duplicate tool pin %q", tool.ID)
		}

		seenToolPins[tool.ID] = true
		if tool.Version == "" {
			return fmt.Errorf("tool pin %q must define version", tool.ID)
		}
	}

	seenRuleIDs := make(map[string]bool, len(config.Rules))
	for _, binding := range config.Rules {
		if binding.RuleID == "" {
			return fmt.Errorf("rule binding has an empty rule_id")
		}

		if seenRuleIDs[binding.RuleID] {
			return fmt.Errorf("duplicate rule binding %q", binding.RuleID)
		}
		seenRuleIDs[binding.RuleID] = true

		switch binding.Level {
		case contract.LevelRequired, contract.LevelRecommendation:
		default:
			return fmt.Errorf("rule %q has invalid level %q", binding.RuleID, binding.Level)
		}

		if !config.Repository.ScopeExists(binding.Scope) {
			return fmt.Errorf("rule %q references unknown scope %q", binding.RuleID, binding.Scope)
		}

		if len(binding.RequirementIDs) == 0 {
			return fmt.Errorf("rule %q must bind at least one requirement", binding.RuleID)
		}

		seenRequirements := make(map[string]bool, len(binding.RequirementIDs))
		for _, requirementID := range binding.RequirementIDs {
			if requirementID == "" {
				return fmt.Errorf("rule %q has an empty requirement ID", binding.RuleID)
			}

			if seenRequirements[requirementID] {
				return fmt.Errorf(
					"rule %q duplicates requirement %q",
					binding.RuleID,
					requirementID,
				)
			}

			seenRequirements[requirementID] = true
		}
	}

	return nil
}

func validateFormatting(formatting policy.FormattingConfig) (err error) {
	headers := formatting.SectionHeaders
	if headers.RequiredMinLines <= 0 {
		return fmt.Errorf("formatting.section_headers.required_min_lines must be positive")
	}

	if headers.ShortFileMaxLines <= 0 {
		return fmt.Errorf("formatting.section_headers.short_file_max_lines must be positive")
	}

	if headers.ShortFileMaxLines >= headers.RequiredMinLines {
		return fmt.Errorf(
			"formatting.section_headers.short_file_max_lines must be less than required_min_lines",
		)
	}

	if headers.OveruseCount <= 0 {
		return fmt.Errorf("formatting.section_headers.overuse_header_count must be positive")
	}

	if len(headers.GenericNames) == 0 {
		return fmt.Errorf("formatting.section_headers.generic_names must not be empty")
	}

	seen := make(map[string]bool, len(headers.GenericNames)+len(headers.StructuralNames))
	for _, names := range [][]string{headers.GenericNames, headers.StructuralNames} {
		for _, name := range names {
			if name == "" {
				return fmt.Errorf("formatting.section_headers contains an empty header name")
			}

			if seen[name] {
				return fmt.Errorf(
					"formatting.section_headers contains duplicate header name %q",
					name,
				)
			}

			seen[name] = true
		}
	}

	return nil
}

func validateGoDomainIdentifiers(identifiers policy.GoDomainIdentifierConfig) (err error) {
	for typeName, constructors := range identifiers {
		if typeName == "" {
			return fmt.Errorf("naming.go_domain_identifiers has an empty type name")
		}

		if len(constructors) == 0 {
			return fmt.Errorf(
				"naming.go_domain_identifiers.%s must define at least one constructor",
				typeName,
			)
		}

		for _, constructor := range constructors {
			if constructor == "" {
				return fmt.Errorf(
					"naming.go_domain_identifiers.%s contains an empty constructor",
					typeName,
				)
			}
		}
	}

	return nil
}
