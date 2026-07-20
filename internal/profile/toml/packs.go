package toml

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
)

const enabledPacksKey = "enabled"

func decodeEnabledPacks(schema map[string]any) (enabledPacks []string, err error) {
	if schema == nil {
		return []string{}, nil
	}

	value, found := schema[enabledPacksKey]
	if !found {
		return []string{}, nil
	}

	enabled, err := decodeStringList(value, "packs.enabled")
	if err != nil {
		return []string{}, err
	}

	return enabled, nil
}

func decodePackConfigs(schema map[string]any) (configs policy.PackConfigs) {
	if schema == nil {
		return nil
	}

	configs = make(policy.PackConfigs, len(schema))
	for packID, config := range schema {
		if packID == enabledPacksKey {
			continue
		}

		section, ok := config.(map[string]any)
		if !ok {
			configs[packID] = policy.PackConfig{"value": config}
			continue
		}

		configs[packID] = policy.PackConfig(section).Clone()
	}

	if len(configs) == 0 {
		return nil
	}

	return configs
}

func encodePacks(
	enabledPacks []string,
	configs policy.PackConfigs,
) (schema map[string]any) {
	if len(enabledPacks) == 0 && configs == nil {
		return nil
	}

	schema = make(map[string]any, len(configs)+1)
	schema[enabledPacksKey] = append([]string{}, enabledPacks...)
	for packID, config := range configs {
		schema[packID] = map[string]any(config.Clone())
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
