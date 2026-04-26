package markers

import "strings"

// Parse parses one inline style-exception marker from a comment line.
func Parse(line string) (marker Marker) {
	directive, found := extractDirective(line)
	if !found {
		return Marker{Status: Absent}
	}

	if !isDirective(directive) {
		return Marker{Status: Invalid}
	}

	body := strings.TrimPrefix(directive, markerPrefix)
	if strings.Contains(body, reasonSeparator) {
		return parseWithReason(body)
	}

	if !isRule(body) {
		return Marker{Status: Invalid}
	}

	return Marker{
		Rule:   body,
		Status: Valid,
	}
}

func parseWithReason(body string) (marker Marker) {
	rule, reason, hasReason := strings.Cut(body, reasonSeparator)
	if !hasReason {
		return Marker{Status: Invalid}
	}

	reason = strings.TrimSpace(reason)
	if !isRule(rule) || reason == "" {
		return Marker{Status: Invalid}
	}

	return Marker{
		Rule:   rule,
		Reason: reason,
		Status: Valid,
	}
}

// Has reports whether the line contains a valid marker for the rule.
func Has(line string, rule string) (valid bool) {
	marker := Parse(line)
	return marker.Status == Valid && marker.Rule == rule
}
