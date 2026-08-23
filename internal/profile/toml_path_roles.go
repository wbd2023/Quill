package profile

func decodePathRoles(schema map[string][]string) (paths PathRoles) {
	return cloneStringLists(PathRoles(schema))
}

func encodePathRoles(paths PathRoles) (schema map[string][]string) {
	return cloneStringLists(map[string][]string(paths))
}
