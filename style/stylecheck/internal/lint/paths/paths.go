package paths

import "strings"

const ClientAdaptersPathSegment = "/internal/client/adapters/"
const RelayAdaptersPathSegment = "/internal/relay/adapters/"
const DomainPathSegment = "/internal/core/domain/"
const DomainPathSuffix = "/internal/core/domain"
const DomainErrorsFilePathSuffix = DomainPathSuffix + "/errors.go"
const ClientServicePathSegment = "/internal/client/application/service/"
const RelayServicePathSegment = "/internal/relay/application/service/"
const ClientPortPathSegment = "/internal/client/application/port/"
const RelayPortPathSegment = "/internal/relay/application/port/"
const CmdPathSegment = "/cmd/"
const InternalPathSegment = "/internal/"
const TestPathSegment = "/test/"
const TestutilMocksPathSegment = "/internal/testkit/mocks/"

func IsAppScopePath(path string) (found bool) {
	return strings.Contains(path, InternalPathSegment) ||
		strings.Contains(path, CmdPathSegment) ||
		strings.Contains(path, TestPathSegment)
}

func IsConcreteInfraPath(path string) (found bool) {
	return strings.Contains(path, ClientAdaptersPathSegment) ||
		strings.Contains(path, RelayAdaptersPathSegment)
}

func IsApplicationPortPath(path string) (found bool) {
	return strings.Contains(path, ClientPortPathSegment) ||
		strings.Contains(path, RelayPortPathSegment)
}

func IsApplicationServicePath(path string) (found bool) {
	return strings.Contains(path, ClientServicePathSegment) ||
		strings.Contains(path, RelayServicePathSegment)
}

func IsDomainPath(path string) (found bool) {
	return strings.Contains(path, DomainPathSegment)
}

func IsDomainErrorsFilePath(path string) (found bool) {
	return strings.HasSuffix(path, DomainErrorsFilePathSuffix)
}

func IsTestMockPath(path string) (found bool) {
	return strings.Contains(path, TestutilMocksPathSegment)
}
