package markers

import "strings"

// marker constants.
const (
	markerPrefix    = "style: "
	reasonSeparator = " because: "
)

// marker constants.
const (
	StatusUnknown Status = ""
	StatusAbsent  Status = "absent"
	StatusInvalid Status = "invalid"
	StatusValid   Status = "valid"
)

// Status is status.
type Status string

// Marker is marker.
type Marker struct {
	Rule   string
	Reason string
	Status Status
}

// Text returns the canonical inline marker for a rule.
func Text(rule string) (marker string) {
	return markerPrefix + rule
}

// Because returns the canonical inline marker with a short justification.
func Because(rule string, reason string) (marker string) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Text(rule)
	}

	return Text(rule) + reasonSeparator + reason
}
