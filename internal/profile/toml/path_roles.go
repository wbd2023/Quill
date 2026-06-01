package toml

import "ciphera/tools/internal/policy"

func decodePathRoles(schema map[string][]string) (paths policy.PathRoles) {
	return cloneStringLists(policy.PathRoles(schema))
}

func encodePathRoles(paths policy.PathRoles) (schema map[string][]string) {
	return cloneStringLists(map[string][]string(paths))
}
