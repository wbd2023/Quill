package profile

import (
	"fmt"
	"strings"

	"github.com/wbd2023/Quill/internal/policy"
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

	if err = validatePathRoles(config.PathRoles); err != nil {
		return err
	}

	if err = validateFileSets(config.Repository, config.FileSets); err != nil {
		return err
	}

	if err = validateTargets(config.Repository, config.Targets); err != nil {
		return err
	}

	if err = validateTools(config.Tools); err != nil {
		return err
	}

	if err = validateEnabledPacks(config.EnabledPacks); err != nil {
		return err
	}

	if err = validatePackConfigs(config.EnabledPacks, config.PackConfigs); err != nil {
		return err
	}

	if err = validateRules(
		config.Repository,
		config.Rules,
	); err != nil {
		return err
	}

	return nil
}

func isBlank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}
