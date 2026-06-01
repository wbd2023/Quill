package drivers

import "strings"

func joinGoLocalImportPrefixes(prefixes []string) (prefix string) {
	return strings.Join(prefixes, ",")
}
