package toml

import (
	"sort"

	"github.com/wbd2023/Quill/internal/style"
)

func decodeScopeMap(source map[string][]string) (target map[style.Scope][]string) {
	if source == nil {
		return nil
	}

	target = make(map[style.Scope][]string, len(source))
	for scope, values := range source {
		target[style.Scope(scope)] = append([]string{}, values...)
	}

	return target
}

func encodeScopeMap(source map[style.Scope][]string) (target map[string][]string) {
	if source == nil {
		return nil
	}

	target = make(map[string][]string, len(source))
	for scope, values := range source {
		target[string(scope)] = append([]string{}, values...)
	}

	return target
}

func sortedMapKeys[V any](source map[string]V) (keys []string) {
	keys = make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneStringLists[M ~map[string][]string](source M) (target M) {
	if source == nil {
		return nil
	}

	target = make(M, len(source))
	for key, values := range source {
		target[key] = append([]string{}, values...)
	}

	return target
}
