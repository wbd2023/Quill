package main

import "strings"

const adaptersPathSegment = "/internal/adapters/"
const domainPathSegment = "/internal/core/domain/"
const domainPathSuffix = "/internal/core/domain"
const domainErrorsFilePathSuffix = domainPathSuffix + "/errors.go"
const corePortsPathSegment = "/internal/core/ports/"
const coreServicesPathSegment = "/internal/core/services/"
const coreServicesAccountRefPathSegment = "/internal/core/services/accountref/"
const cmdPathSegment = "/cmd/"
const internalPathSegment = "/internal/"
const mocksPathSegment = "/internal/mocks/"
const testsPathSegment = "/tests/"

func isAppScopePath(path string) (found bool) {
	return strings.Contains(path, internalPathSegment) ||
		strings.Contains(path, cmdPathSegment) ||
		strings.Contains(path, testsPathSegment)
}
