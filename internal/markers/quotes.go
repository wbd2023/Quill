package markers

import "strings"

// quotes constants.
const (
	escapeCharacter = '\\'
	doubleQuote     = '"'
	singleQuote     = '\''
	backtickQuote   = '`'
)

func indexOutsideQuotes(line string, needle string) (index int) {
	index = 0
	for index < len(line) {
		next := strings.Index(line[index:], needle)
		if next < 0 {
			return -1
		}
		next += index

		if !insideQuotedText(line[:next]) {
			return next
		}

		index = next + len(needle)
	}

	return -1
}

func insideQuotedText(prefix string) (inside bool) {
	var quote rune
	escaped := false

	for _, character := range prefix {
		if escaped {
			escaped = false
			continue
		}

		if quote != 0 {
			if quote != backtickQuote && character == escapeCharacter {
				escaped = true
				continue
			}

			if character == quote {
				quote = 0
			}

			continue
		}

		if isQuote(character) {
			quote = character
		}
	}

	return quote != 0
}

func isQuote(character rune) (quote bool) {
	return character == doubleQuote || character == singleQuote || character == backtickQuote
}
