package profile

import "fmt"

/* ---------------------------------------- Pack Policies --------------------------------------- */

// PackPolicies stores raw profile policy owned by Packs.
type PackPolicies map[string]PackPolicy

// PackPolicy stores one Pack's raw profile policy object.
type PackPolicy map[string]any

// Lookup returns the named Pack Policy.
func (policies PackPolicies) Lookup(packID string) (policy PackPolicy, found bool) {
	policy, found = policies[packID]
	return policy, found
}

// Clone returns a deep copy of policies.
func (policies PackPolicies) Clone() (clone PackPolicies) {
	if policies == nil {
		return nil
	}

	clone = make(PackPolicies, len(policies))
	for packID, policy := range policies {
		clone[packID] = policy.Clone()
	}

	return clone
}

// Clone returns a deep copy of policy.
func (policy PackPolicy) Clone() (clone PackPolicy) {
	if policy == nil {
		return nil
	}

	clone = make(PackPolicy, len(policy))
	for key, value := range policy {
		clone[key] = clonePackPolicyValue(value)
	}

	return clone
}

/* --------------------------------------- Policy Cloning --------------------------------------- */

func clonePackPolicyValue(value any) (clone any) {
	switch value := value.(type) {
	case PackPolicy:
		return value.Clone()

	case map[string]any:
		clone := make(map[string]any, len(value))
		for key, child := range value {
			clone[key] = clonePackPolicyValue(child)
		}
		return clone

	case []map[string]any:
		clone := make([]map[string]any, 0, len(value))
		for _, child := range value {
			clone = append(clone, map[string]any(PackPolicy(child).Clone()))
		}
		return clone

	case []any:
		clone := make([]any, 0, len(value))
		for _, child := range value {
			clone = append(clone, clonePackPolicyValue(child))
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

/* ----------------------------------------- Validation ----------------------------------------- */

func validatePackPolicies(
	enabledPacks []string,
	policies PackPolicies,
) (err error) {
	if len(policies) == 0 {
		return nil
	}

	enabled := make(map[string]bool, len(enabledPacks))
	for _, packID := range enabledPacks {
		enabled[packID] = true
	}

	for packID := range policies {
		if isBlank(packID) {
			return fmt.Errorf("packs contains an empty pack id")
		}

		if !enabled[packID] {
			return fmt.Errorf("packs.%s policy is not enabled in packs.enabled", packID)
		}

	}

	return nil
}
