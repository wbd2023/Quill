package gopolicy

import "fmt"

func stringList(section map[string]any, key string, field string) (values []string, err error) {
	if section == nil {
		return nil, nil
	}

	value, found := section[key]
	if !found {
		return nil, nil
	}

	return decodeStringList(value, field)
}

func stringListMap(
	section map[string]any,
	key string,
	field string,
) (values DomainValueConstructors, err error) {
	if section == nil {
		return nil, nil
	}

	value, found := section[key]
	if !found {
		return nil, nil
	}

	raw, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a table", field)
	}

	values = make(DomainValueConstructors, len(raw))
	for name, constructors := range raw {
		values[name], err = decodeStringList(constructors, field+"."+name)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func decodeStringList(value any, field string) (values []string, err error) {
	switch items := value.(type) {
	case []any:
		values = make([]string, 0, len(items))
		for _, item := range items {
			text, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s must contain only strings", field)
			}

			values = append(values, text)
		}

		return values, nil

	case []string:
		return cloneStrings(items), nil

	default:
		return nil, fmt.Errorf("%s must be a string array", field)
	}
}

func encodeStringListMap(values DomainValueConstructors) (encoded map[string]any) {
	if values == nil {
		return nil
	}

	encoded = make(map[string]any, len(values))
	for name, constructors := range values {
		encoded[name] = cloneStrings(constructors)
	}

	return encoded
}

func cloneStrings(values []string) (clone []string) {
	return append([]string{}, values...)
}
