package markers

import "strings"

const (
	rulePrefix       = "allow-"
	asciiRuneMaximum = 127
)

/* ------------------------------------------- Parsing ------------------------------------------ */

// Parse parses one inline style-exception marker from a comment line.
//
// Marker grammar: style: allow-<rule> [because: <reason>].
func Parse(line string) (marker Marker) {
	directive, found := extractDirective(line)
	if !found {
		return Marker{Status: StatusAbsent}
	}

	if !isMarkerDirective(directive) {
		return Marker{Status: StatusInvalid}
	}

	return parseBody(strings.TrimPrefix(directive, markerPrefix))
}

// Has reports whether the line contains a valid marker for the rule.
func Has(line string, rule string) (valid bool) {
	marker := Parse(line)
	return marker.Status == StatusValid && marker.Rule == rule
}

func parseBody(body string) (marker Marker) {
	rule, reason, hasReason := strings.Cut(body, reasonSeparator)
	if !isRule(rule) {
		return Marker{Status: StatusInvalid}
	}

	if !hasReason {
		return Marker{Rule: rule, Status: StatusValid}
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Marker{Status: StatusInvalid}
	}

	return Marker{Rule: rule, Reason: reason, Status: StatusValid}
}

/* ----------------------------------------- Validation ----------------------------------------- */

func isMarkerDirective(value string) (valid bool) {
	return strings.HasPrefix(value, markerPrefix) && isASCIIText(value)
}

func isRule(value string) (valid bool) {
	if !strings.HasPrefix(value, rulePrefix) {
		return false
	}

	name := value[len(rulePrefix):]
	for part := range strings.SplitSeq(name, "-") {
		if !isRulePart(part) {
			return false
		}
	}

	return true
}

func isRulePart(value string) (valid bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if isASCIILower(character) || isASCIIDigit(character) {
			continue
		}

		return false
	}

	return true
}

/* ----------------------------------------- Characters ----------------------------------------- */

func isASCIILower(character rune) (lower bool) {
	return 'a' <= character && character <= 'z'
}

func isASCIIDigit(character rune) (digit bool) {
	return '0' <= character && character <= '9'
}

func isASCIIText(value string) (ascii bool) {
	for _, character := range value {
		if character > asciiRuneMaximum {
			return false
		}
	}

	return true
}
