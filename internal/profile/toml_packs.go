package profile

import "fmt"

// EnabledPacksKey is reserved in the [packs] table for the selected Pack IDs.
const EnabledPacksKey = "enabled"

func decodeEnabledPacks(schema map[string]any) (enabledPacks []string, err error) {
	if schema == nil {
		return []string{}, nil
	}

	value, found := schema[EnabledPacksKey]
	if !found {
		return []string{}, nil
	}

	enabled, err := decodeStringList(value, "packs.enabled")
	if err != nil {
		return []string{}, err
	}

	return enabled, nil
}

func decodePackPolicies(schema map[string]any) (policies PackPolicies, err error) {
	if schema == nil {
		return nil, nil
	}

	policies = make(PackPolicies, len(schema))
	for packID, policy := range schema {
		if packID == EnabledPacksKey {
			continue
		}

		section, ok := policy.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("packs.%s policy must be a table", packID)
		}

		policies[packID] = PackPolicy(section).Clone()
	}

	if len(policies) == 0 {
		return nil, nil
	}

	return policies, nil
}

func encodePacks(
	enabledPacks []string,
	policies PackPolicies,
) (schema map[string]any) {
	if len(enabledPacks) == 0 && policies == nil {
		return nil
	}

	schema = make(map[string]any, len(policies)+1)
	schema[EnabledPacksKey] = append([]string{}, enabledPacks...)
	for packID, policy := range policies {
		schema[packID] = map[string]any(policy.Clone())
	}

	return schema
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
		return append([]string{}, items...), nil

	default:
		return nil, fmt.Errorf("%s must be a string array", field)
	}
}
