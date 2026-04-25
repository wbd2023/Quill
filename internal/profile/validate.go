package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func (policy Profile) Validate() (err error) {
	if policy.SchemaVersion != SchemaVersion1 {
		return fmt.Errorf("unsupported style profile version %d", policy.SchemaVersion)
	}

	if len(policy.RulePacks.Enabled) == 0 {
		return fmt.Errorf("rule_packs.enabled must not be empty")
	}

	if len(policy.Repository.RootMarkers) == 0 {
		return fmt.Errorf("repository.root_markers must not be empty")
	}

	if len(policy.Repository.AppScanRoots) == 0 {
		return fmt.Errorf("repository.app_scan_roots must not be empty")
	}

	if len(policy.Repository.ToolsScanRoots) == 0 {
		return fmt.Errorf("repository.tools_scan_roots must not be empty")
	}

	if policy.Repository.GeneratedMarker == "" {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	if policy.Repository.GeneratedProbeLimit <= 0 {
		return fmt.Errorf("repository.generated_probe_limit must be positive")
	}

	if policy.StyleGuide.Path == "" {
		return fmt.Errorf("styleguide.path must not be empty")
	}

	if policy.StyleGuide.RequirementIDFormat == "" {
		return fmt.Errorf("styleguide.requirement_id_format must not be empty")
	}

	if policy.StyleGuide.RequirementIDFormat != RequirementIDFormatSectionSlug {
		return fmt.Errorf(
			"unsupported styleguide.requirement_id_format %q",
			policy.StyleGuide.RequirementIDFormat,
		)
	}

	if policy.Imports.LocalPrefix == "" {
		return fmt.Errorf("imports.local_prefix must not be empty")
	}

	if err = validateGoDomainIdentifiers(policy.Naming.GoDomainIdentifiers); err != nil {
		return err
	}

	seenFileSets := make(map[string]bool, len(policy.FileSets))
	for _, fileSet := range policy.FileSets {
		if fileSet.Name == "" {
			return fmt.Errorf("file set name must not be empty")
		}

		if seenFileSets[fileSet.Name] {
			return fmt.Errorf("duplicate file set %q", fileSet.Name)
		}

		seenFileSets[fileSet.Name] = true
	}

	seenLanguageBackends := make(map[string]bool, len(policy.Language.Backends))
	for _, backend := range policy.Language.Backends {
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
	}

	if policy.ControlPlane.QualityFile == "" {
		return fmt.Errorf("control_plane.quality_file must not be empty")
	}

	if len(policy.Architecture.Layers) == 0 {
		return fmt.Errorf("architecture.layers must not be empty")
	}

	if len(policy.Rules) == 0 {
		return fmt.Errorf("rules must not be empty")
	}

	layerNames := make(map[string]bool, len(policy.Architecture.Layers))
	for _, layer := range policy.Architecture.Layers {
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

	for _, layer := range policy.Architecture.Layers {
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

	seenRuleIDs := make(map[string]bool, len(policy.Rules))
	for _, binding := range policy.Rules {
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

		switch binding.Scope {
		case contract.ScopeApp, contract.ScopeTools, contract.ScopeAll:
		default:
			return fmt.Errorf("rule %q has invalid scope %q", binding.RuleID, binding.Scope)
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

func validateGoDomainIdentifiers(identifiers GoDomainIdentifierConfig) (err error) {
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
