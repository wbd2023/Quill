package markers

import "strings"

const (
	LongLine = "allow-long-line"
	NonASCII = "allow-non-ascii"
)

const (
	Absent  Status = "absent"
	Invalid Status = "invalid"
	Valid   Status = "valid"
)

type Status string

type Marker struct {
	Rule   string
	Reason string
	Status Status
}

// Text returns the canonical inline marker for a rule.
func Text(rule string) (marker string) {
	return markerPrefix + rule
}

// TextWithReason returns the canonical marker with a short justification.
func TextWithReason(rule string, reason string) (marker string) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Text(rule)
	}

	return Text(rule) + reasonSeparator + reason
}
