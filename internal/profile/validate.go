package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

// Validate checks config for supported schema version and internal consistency.
func Validate(config policy.Config) (err error) {
	if config.SchemaVersion != policy.SchemaVersion {
		return fmt.Errorf("unsupported style profile version %d", config.SchemaVersion)
	}

	if err = validateRepository(config.Repository); err != nil {
		return err
	}

	if err = validateStyleGuide(config.StyleGuide); err != nil {
		return err
	}

	if err = validateFileSets(config.Repository, config.FileSets); err != nil {
		return err
	}

	if err = validateLanguage(config.Repository, config.Language); err != nil {
		return err
	}

	if err = validateGo(config.Go, config.Language); err != nil {
		return err
	}

	if err = validateTools(config.Tools); err != nil {
		return err
	}

	if err = validateFormatting(config.Formatting); err != nil {
		return err
	}

	if err = validateVocabulary(config.Vocabulary); err != nil {
		return err
	}

	if err = validateQualitySurface(config.QualitySurface); err != nil {
		return err
	}

	if err = validateRulePacks(config.RulePacks); err != nil {
		return err
	}

	if err = validateRules(config.Repository, config.Rules); err != nil {
		return err
	}

	return nil
}
