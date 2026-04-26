package markers

import "strings"

func indexOutsideQuotedText(line string, needle string) (position int) {
	position = 0
	for position < len(line) {
		next := strings.Index(line[position:], needle)
		if next < 0 {
			return -1
		}
		next += position

		if !insideQuotedText(line[:next]) {
			return next
		}

		position = next + len(needle)
	}

	return -1
}

func insideQuotedText(prefix string) (inside bool) {
	escaped := false
	for _, character := range prefix {
		if escaped {
			escaped = false
			continue
		}

		if character == '\\' {
			escaped = true
			continue
		}

		if character == '"' {
			inside = !inside
		}
	}

	return inside
}
