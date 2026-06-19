package vocabularypolicy

import "fmt"

func cloneStrings(values []string) (clone []string) {
	return append([]string{}, values...)
}

func validateList(field string, values []string) (err error) {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if value == "" {
			return fmt.Errorf("%s contains an empty value", field)
		}

		if seen[value] {
			return fmt.Errorf("%s contains duplicate value %q", field, value)
		}

		seen[value] = true
	}

	return nil
}

func stringListMap(
	section map[string]any,
	key string,
	field string,
) (values map[string][]string, err error) {
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

	values = make(map[string][]string, len(raw))
	for name, forbidden := range raw {
		values[name], err = decodeStringList(forbidden, field+"."+name)
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

func encodeStringListMap(values map[string][]string) (encoded map[string]any) {
	if values == nil {
		return nil
	}

	encoded = make(map[string]any, len(values))
	for name, forbidden := range values {
		encoded[name] = cloneStrings(forbidden)
	}

	return encoded
}
