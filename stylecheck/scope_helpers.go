package main

import "strings"

const cmdPathSegment = "/cmd/"
const internalPathSegment = "/internal/"
const testsPathSegment = "/tests/"

func isAppScopePath(path string) (found bool) {
	return strings.Contains(path, internalPathSegment) ||
		strings.Contains(path, cmdPathSegment) ||
		strings.Contains(path, testsPathSegment)
}
