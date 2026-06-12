package gopolicy

import "fmt"

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

func boolField(section map[string]any, key string, field string) (enabled bool, err error) {
	if section == nil {
		return false, nil
	}

	value, found := section[key]
	if !found {
		return false, nil
	}

	enabled, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", field)
	}

	return enabled, nil
}
