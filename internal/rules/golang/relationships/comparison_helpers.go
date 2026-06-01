package relationships

func methodNames(methods []methodDeclaration) (names []string) {
	names = make([]string, len(methods))
	for index, method := range methods {
		names[index] = method.Name
	}

	return names
}

func methodNameSet(names []string) (set map[string]bool) {
	set = make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}

	return set
}

func matchingMethods(
	methods []methodDeclaration,
	names map[string]bool,
) (matches []methodDeclaration) {
	matches = make([]methodDeclaration, 0, len(names))
	for _, method := range methods {
		if names[method.Name] {
			matches = append(matches, method)
		}
	}

	return matches
}
