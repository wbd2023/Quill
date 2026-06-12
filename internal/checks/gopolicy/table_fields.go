package gopolicy

import "fmt"

func configSection(
	pack map[string]any,
	key string,
	field string,
) (section map[string]any, err error) {
	value, found := pack[key]
	if !found {
		return nil, nil
	}

	section, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a table", field)
	}

	return section, nil
}

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
