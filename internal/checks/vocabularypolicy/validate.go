package vocabularypolicy

import "fmt"

// ValidateConfig validates Vocabulary Pack Policy.
func ValidateConfig(config Config) (err error) {
	if err = validateStringListMap(
		"packs.vocabulary.go.type_suffixes",
		config.Go.TypeSuffixes,
	); err != nil {
		return err
	}

	if err = validateStringListMap(
		"packs.vocabulary.go.identifier_suffixes",
		config.Go.IdentifierSuffixes,
	); err != nil {
		return err
	}

	if err = validateStringListMap(
		"packs.vocabulary.bash.variable_names",
		config.Bash.VariableNames,
	); err != nil {
		return err
	}

	return nil
}

func validateStringListMap(field string, values map[string][]string) (err error) {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(values))
	for preferred, forbidden := range values {
		if preferred == "" {
			return fmt.Errorf("%s has an empty preferred name", field)
		}

		if err = validateList(field+"."+preferred, forbidden); err != nil {
			return err
		}

		for _, shorthand := range forbidden {
			if seen[shorthand] {
				return fmt.Errorf(
					"%s maps shorthand %q to more than one preferred name",
					field,
					shorthand,
				)
			}

			seen[shorthand] = true
		}
	}

	return nil
}
