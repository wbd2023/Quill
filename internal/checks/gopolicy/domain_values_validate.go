package gopolicy

import "fmt"

func validateDomainValueConstructors(
	constructors DomainValueConstructors,
) (err error) {
	for typeName, typeConstructors := range constructors {
		if blank(typeName) {
			return fmt.Errorf("packs.go.domain_values.required_constructors has an empty type name")
		}

		if len(typeConstructors) == 0 {
			field := "packs.go.domain_values.required_constructors." + typeName
			return fmt.Errorf(
				"%s must define at least one constructor",
				field,
			)
		}

		if err = validateList(
			"packs.go.domain_values.required_constructors."+typeName,
			typeConstructors,
		); err != nil {
			return err
		}
	}

	return nil
}
