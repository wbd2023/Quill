package policy

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
