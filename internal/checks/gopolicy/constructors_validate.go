package gopolicy

import "fmt"

func validateConstructors(config ConstructorConfig) (err error) {
	seen := make(map[string]bool, len(config.ParameterOrder))
	for _, group := range config.ParameterOrder {
		if blank(group.Name) {
			return fmt.Errorf("packs.go.constructors.parameter_order contains an empty name")
		}

		if seen[group.Name] {
			return fmt.Errorf(
				"packs.go.constructors.parameter_order contains duplicate name %q",
				group.Name,
			)
		}

		seen[group.Name] = true

		if !group.MatchesSecretNames &&
			len(group.TypeNameSuffixes) == 0 &&
			len(group.ParameterNames) == 0 {
			return fmt.Errorf(
				"packs.go.constructors.parameter_order.%s must define at least one matcher",
				group.Name,
			)
		}

		if err = validateList(
			"packs.go.constructors.parameter_order."+group.Name+".type_name_suffixes",
			group.TypeNameSuffixes,
		); err != nil {
			return err
		}

		if err = validateList(
			"packs.go.constructors.parameter_order."+group.Name+".parameter_names",
			group.ParameterNames,
		); err != nil {
			return err
		}
	}

	return nil
}
