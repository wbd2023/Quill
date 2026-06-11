package policy

import "fmt"

func stringList(section map[string]any, key string, field string) (values []string, err error) {
	if section == nil {
		return nil, nil
	}

	value, found := section[key]
	if !found {
		return nil, nil
	}

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

func stringField(section map[string]any, key string, field string) (text string, err error) {
	if section == nil {
		return "", nil
	}

	value, found := section[key]
	if !found {
		return "", nil
	}

	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", field)
	}

	return text, nil
}
