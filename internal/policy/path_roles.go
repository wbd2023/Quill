package policy

// PathRoles maps profile-owned role names to path patterns.
type PathRoles map[string][]string

// LookupPatterns returns a defensive copy of the patterns for the named path role.
func (p PathRoles) LookupPatterns(name string) (patterns []string) {
	values, found := p[name]
	if !found {
		return nil
	}

	return append([]string{}, values...)
}
