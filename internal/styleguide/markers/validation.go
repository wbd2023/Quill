package markers

import "strings"

const (
	rulePrefix       = "allow-"
	ruleSeparator    = '-'
	asciiRuneMaximum = 127
)

func isDirective(value string) (valid bool) {
	return strings.HasPrefix(value, markerPrefix) && isASCII(value)
}

func isRule(value string) (valid bool) {
	if !strings.HasPrefix(value, rulePrefix) {
		return false
	}

	for _, character := range value[len(rulePrefix):] {
		if isRuleCharacter(character) {
			continue
		}

		return false
	}

	return len(value) > len(rulePrefix)
}

func isRuleCharacter(character rune) (valid bool) {
	return character == ruleSeparator || isLowerASCII(character) || isDigit(character)
}

func isLowerASCII(character rune) (lower bool) {
	return 'a' <= character && character <= 'z'
}

func isDigit(character rune) (digit bool) {
	return '0' <= character && character <= '9'
}

func isASCII(value string) (ascii bool) {
	for _, character := range value {
		if character > asciiRuneMaximum {
			return false
		}
	}

	return true
}
