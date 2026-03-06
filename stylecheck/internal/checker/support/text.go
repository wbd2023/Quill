package support

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const sentenceEndingPunctuation = ".!?"

func StartsWithUppercaseLetter(value string) (found bool) {
	firstRune, _ := utf8.DecodeRuneInString(value)
	return unicode.IsUpper(firstRune)
}

func EndsWithSentencePunctuation(value string) (found bool) {
	lastRune, _ := utf8.DecodeLastRuneInString(value)
	return strings.ContainsRune(sentenceEndingPunctuation, lastRune)
}
