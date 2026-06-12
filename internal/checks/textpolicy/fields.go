package textpolicy

import (
	"fmt"
	"slices"
)

func rejectUnknownFields(
	section map[string]any,
	field string,
	allowed ...string,
) (err error) {
	known := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		known[name] = true
	}

	var unknown []string
	for name := range section {
		if !known[name] {
			unknown = append(unknown, name)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	slices.Sort(unknown)
	return fmt.Errorf("%s.%s is not supported", field, unknown[0])
}

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

func intField(section map[string]any, key string, field string) (number int, err error) {
	if section == nil {
		return 0, nil
	}

	value, found := section[key]
	if !found {
		return 0, nil
	}

	switch value := value.(type) {
	case int:
		return value, nil

	case int64:
		number = int(value)
		if int64(number) != value {
			return 0, fmt.Errorf("%s is outside the supported integer range", field)
		}

		return number, nil

	default:
		return 0, fmt.Errorf("%s must be an integer", field)
	}
}

func cloneStrings(values []string) (clone []string) {
	return append([]string{}, values...)
}
