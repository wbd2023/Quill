package policy

type PathClassSet struct {
	Classes map[string][]string
}

func (paths PathClassSet) Patterns(className string) (patterns []string) {
	if paths.Classes == nil {
		return nil
	}

	return append([]string{}, paths.Classes[className]...)
}
