package project

import (
	"fmt"
	"slices"
)

/* --------------------------------------- Unknown Fields --------------------------------------- */

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

/* ------------------------------------------ Sections ------------------------------------------ */

func configSection(
	section map[string]any,
	key string,
	field string,
) (child map[string]any, err error) {
	value, found := section[key]
	if !found {
		return nil, nil
	}

	child, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a table", field)
	}

	return child, nil
}

/* ------------------------------------------- Scalars ------------------------------------------ */

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

/* ------------------------------------------- Tables ------------------------------------------- */

func tableList(
	section map[string]any,
	key string,
	field string,
) (tables []map[string]any, err error) {
	if section == nil {
		return nil, nil
	}

	value, found := section[key]
	if !found {
		return nil, nil
	}

	switch items := value.(type) {
	case []map[string]any:
		return append([]map[string]any{}, items...), nil

	case []any:
		tables = make([]map[string]any, 0, len(items))
		for _, item := range items {
			table, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%s must contain only tables", field)
			}

			tables = append(tables, table)
		}

		return tables, nil

	default:
		return nil, fmt.Errorf("%s must be an array of tables", field)
	}
}
