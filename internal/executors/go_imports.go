package executors

import "strings"

func joinGoLocalImportPrefixes(prefixes []string) (prefix string) {
	return strings.Join(prefixes, ",")
}
