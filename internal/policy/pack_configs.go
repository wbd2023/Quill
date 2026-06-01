package policy

// PackConfigs stores raw profile config owned by Packs.
type PackConfigs map[string]PackConfig

// PackConfig stores one Pack's raw profile config subtree.
type PackConfig map[string]any

// Lookup returns the named Pack config.
func (configs PackConfigs) Lookup(packID string) (config PackConfig, found bool) {
	config, found = configs[packID]
	return config, found
}

// Clone returns a deep copy of configs.
func (configs PackConfigs) Clone() (clone PackConfigs) {
	if configs == nil {
		return nil
	}

	clone = make(PackConfigs, len(configs))
	for packID, config := range configs {
		clone[packID] = config.Clone()
	}

	return clone
}

// Clone returns a deep copy of config.
func (config PackConfig) Clone() (clone PackConfig) {
	if config == nil {
		return nil
	}

	clone = make(PackConfig, len(config))
	for key, value := range config {
		clone[key] = clonePackConfigValue(value)
	}

	return clone
}

func clonePackConfigValue(value any) (clone any) {
	switch value := value.(type) {
	case PackConfig:
		return value.Clone()

	case map[string]any:
		clone := make(map[string]any, len(value))
		for key, child := range value {
			clone[key] = clonePackConfigValue(child)
		}
		return clone

	case []map[string]any:
		clone := make([]map[string]any, 0, len(value))
		for _, child := range value {
			clone = append(clone, map[string]any(PackConfig(child).Clone()))
		}
		return clone

	case []any:
		clone := make([]any, 0, len(value))
		for _, child := range value {
			clone = append(clone, clonePackConfigValue(child))
		}
		return clone

	case []string:
		return append([]string{}, value...)

	case []int64:
		return append([]int64{}, value...)

	case []int:
		return append([]int{}, value...)

	case []float64:
		return append([]float64{}, value...)

	case []bool:
		return append([]bool{}, value...)

	default:
		return value
	}
}
