package styleguide

import "strings"

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	ExceptionLongLine = "allow-long-line"
	ExceptionNonASCII = "allow-non-ascii"
	asciiMaximum      = 127
)

const (
	exceptionPrefix          = "style: "
	exceptionReasonSeparator = " because: "
)

/* -------------------------------------- Exception Markers ------------------------------------- */

// ExceptionMarker returns the canonical inline style-exception marker for a rule.
func ExceptionMarker(rule string) (marker string) {
	return exceptionPrefix + rule
}

// ExceptionMarkerWithReason returns the canonical marker with a short justification.
func ExceptionMarkerWithReason(rule string, reason string) (marker string) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return ExceptionMarker(rule)
	}

	return ExceptionMarker(rule) + exceptionReasonSeparator + reason
}

/* -------------------------------------- Exception Parsing ------------------------------------- */

// ParseExceptionMarker parses one inline style-exception marker from a comment line.
func ParseExceptionMarker(line string) (rule string, reason string, found bool, valid bool) {
	directive, found := extractExceptionDirective(line)
	if !found {
		return "", "", false, false
	}

	if !strings.HasPrefix(directive, exceptionPrefix) {
		return "", "", true, false
	}

	if !isASCII(directive) {
		return "", "", true, false
	}

	body := strings.TrimPrefix(directive, exceptionPrefix)
	if strings.Contains(body, exceptionReasonSeparator) {
		exceptionRule, reason, hasReason := strings.Cut(body, exceptionReasonSeparator)
		if !hasReason {
			return "", "", true, false
		}

		reason = strings.TrimSpace(reason)
		if !isExceptionRule(exceptionRule) || reason == "" {
			return "", "", true, false
		}

		return exceptionRule, reason, true, true
	}

	if !isExceptionRule(body) {
		return "", "", true, false
	}

	return body, "", true, true
}

/* --------------------------------------- Marker Queries --------------------------------------- */

// HasExceptionMarker reports whether the line contains a valid marker for the rule.
func HasExceptionMarker(line string, rule string) (valid bool) {
	exceptionRule, _, found, ok := ParseExceptionMarker(line)
	return found && ok && exceptionRule == rule
}

/* ---------------------------------------- ASCII Checks ---------------------------------------- */

func isASCII(value string) (ascii bool) {
	for _, runeValue := range value {
		if runeValue > asciiMaximum {
			return false
		}
	}

	return true
}

/* ------------------------------------ Directive Extraction ------------------------------------ */

func extractExceptionDirective(line string) (directive string, found bool) {
	for _, prefix := range []string{
		"# " + exceptionPrefix,
		"// " + exceptionPrefix,
		"/* " + exceptionPrefix,
		"* " + exceptionPrefix,
	} {
		index := indexOutsideQuotedText(line, prefix)
		if index < 0 {
			continue
		}

		directive = strings.TrimSpace(line[index+len(prefix)-len(exceptionPrefix):])
		directive = strings.TrimSpace(strings.TrimSuffix(directive, "*/"))
		return directive, true
	}

	return "", false
}

/* --------------------------------------- Rule Validation -------------------------------------- */

func isExceptionRule(value string) (valid bool) {
	if !strings.HasPrefix(value, "allow-") {
		return false
	}

	for _, runeValue := range value[len("allow-"):] {
		if runeValue == '-' ||
			('a' <= runeValue && runeValue <= 'z') ||
			('0' <= runeValue && runeValue <= '9') {
			continue
		}

		return false
	}

	return len(value) > len("allow-")
}

/* ------------------------------------ Quoted Text Scanning ------------------------------------ */

func indexOutsideQuotedText(line string, token string) (index int) {
	const noEscape = byte(0)

	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	escapePrefix := noEscape

	for current := 0; current < len(line); current++ {
		character := line[current]
		switch {
		case escapePrefix != noEscape:
			escapePrefix = noEscape
			continue

		case inSingleQuote:
			if character == '\'' {
				inSingleQuote = false
			}
			continue

		case inDoubleQuote:
			if character == '\\' {
				escapePrefix = character
				continue
			}

			if character == '"' {
				inDoubleQuote = false
			}
			continue

		case inBacktick:
			if character == '`' {
				inBacktick = false
			}
			continue

		case character == '\'':
			inSingleQuote = true
			continue

		case character == '"':
			inDoubleQuote = true
			continue

		case character == '`':
			inBacktick = true
			continue
		}

		if strings.HasPrefix(line[current:], token) {
			return current
		}
	}

	return -1
}
