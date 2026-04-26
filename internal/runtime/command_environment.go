package runtime

import "sort"

func environmentEntries(environment map[string]string) (values []string) {
	if len(environment) == 0 {
		return nil
	}

	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values = make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+environment[key])
	}

	return values
}
