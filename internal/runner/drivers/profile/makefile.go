package profile

import "strings"

type makefileSurface struct {
	Variables map[string]string
	Targets   map[string]makefileTarget
}

type makefileTarget struct {
	Recipes []string
}

func parseMakefileSurface(contents string) (surface makefileSurface) {
	surface = makefileSurface{
		Variables: make(map[string]string),
		Targets:   make(map[string]makefileTarget),
	}

	activeTarget := ""
	for _, line := range strings.Split(contents, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "", strings.HasPrefix(trimmed, "#"):
			activeTarget = ""

		case strings.HasPrefix(line, "\t"):
			if activeTarget == "" {
				continue
			}

			target := surface.Targets[activeTarget]
			target.Recipes = append(target.Recipes, strings.TrimSpace(line))
			surface.Targets[activeTarget] = target

		case strings.Contains(trimmed, "=") && !strings.Contains(trimmed, ":"):
			name, value, _ := strings.Cut(trimmed, "=")
			surface.Variables[strings.TrimSpace(name)] = strings.TrimSpace(value)
			activeTarget = ""

		case strings.Contains(trimmed, ":"):
			targetName, _, _ := strings.Cut(trimmed, ":")
			targetName = strings.TrimSpace(targetName)
			surface.Targets[targetName] = makefileTarget{}
			activeTarget = targetName

		default:
			activeTarget = ""
		}
	}

	return surface
}

func hasRecipeLine(lines []string, expected string) (found bool) {
	for _, line := range lines {
		if strings.TrimSpace(line) == expected {
			return true
		}
	}

	return false
}
