package policy

import (
	"fmt"
	"strings"
)

// ValidateConfig validates Go rule policy.
func ValidateConfig(config Config) (err error) {
	if len(config.LocalImportPrefixes) == 0 {
		return fmt.Errorf("packs.go.local_import_prefixes must not be empty")
	}

	if err = validateList(
		"packs.go.local_import_prefixes",
		config.LocalImportPrefixes,
	); err != nil {
		return err
	}

	if err = validateParameters(config.Parameters); err != nil {
		return err
	}

	if err = validateConstructors(config.Constructors); err != nil {
		return err
	}

	if err = validateDomainValueConstructors(
		config.DomainValues.RequiredConstructors,
	); err != nil {
		return err
	}

	return validateArchitecture(config.Architecture)
}

func validateList(field string, values []string) (err error) {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if blank(value) {
			return fmt.Errorf("%s contains an empty value", field)
		}

		if seen[value] {
			return fmt.Errorf("%s contains duplicate value %q", field, value)
		}

		seen[value] = true
	}

	return nil
}

func blank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}
