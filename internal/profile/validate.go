package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func Validate(config policy.Config) (err error) {
	if config.SchemaVersion != policy.SchemaVersion {
		return fmt.Errorf("unsupported style profile version %d", config.SchemaVersion)
	}

	if len(config.RulePacks.Enabled) == 0 {
		return fmt.Errorf("rule_packs.enabled must not be empty")
	}

	if err = validateRepository(config.Repository); err != nil {
		return err
	}

	if err = validateStyleGuide(config.StyleGuide); err != nil {
		return err
	}

	if err = validateFormatting(config.Formatting); err != nil {
		return err
	}

	if err = validateImports(config.Imports); err != nil {
		return err
	}

	if err = validateFileSets(config.Repository, config.FileSets); err != nil {
		return err
	}

	if err = validateLanguage(config.Repository, config.Language); err != nil {
		return err
	}

	if err = validateTools(config.Tools); err != nil {
		return err
	}

	if err = validateNaming(config.Naming); err != nil {
		return err
	}

	if err = validateControlPlane(config.ControlPlane); err != nil {
		return err
	}

	if err = validateArchitecture(config.Architecture); err != nil {
		return err
	}

	if err = validateRules(config.Repository, config.Rules); err != nil {
		return err
	}

	return nil
}
