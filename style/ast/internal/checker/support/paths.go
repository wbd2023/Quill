package support

import "strings"

const AdaptersPathSegment = "/internal/adapters/"
const DomainPathSegment = "/internal/core/domain/"
const DomainPathSuffix = "/internal/core/domain"
const DomainErrorsFilePathSuffix = DomainPathSuffix + "/errors.go"
const CorePortsPathSegment = "/internal/core/ports/"
const CoreServicesPathSegment = "/internal/core/services/"
const CoreServicesAccountRefPathSegment = "/internal/core/services/accountref/"
const cmdPathSegment = "/cmd/"
const internalPathSegment = "/internal/"
const testsPathSegment = "/tests/"
const MocksPathSegment = "/internal/mocks/"

func IsAppScopePath(path string) (found bool) {
	return strings.Contains(path, internalPathSegment) ||
		strings.Contains(path, cmdPathSegment) ||
		strings.Contains(path, testsPathSegment)
}
