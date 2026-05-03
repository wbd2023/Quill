package policy

// PathClasses maps profile-owned class names to path patterns.
type PathClasses map[string][]string

// LookupPatterns returns a defensive copy of the patterns for the named path class.
func (p PathClasses) LookupPatterns(name string) (patterns []string) {
	values, found := p[name]
	if !found {
		return nil
	}

	return append([]string{}, values...)
}
