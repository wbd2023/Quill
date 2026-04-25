package rulepack

/* --------------------------------------- Go Path Classes -------------------------------------- */

const (
	PathClassApp             = "app"
	PathClassApplicationPort = "application_port"
	PathClassConcreteInfra   = "concrete_infra"
	PathClassDomain          = "domain"
	PathClassDomainErrors    = "domain_errors"
	PathClassTestMocks       = "test_mocks"
)

func requiredGoStylePathClasses() (classes []string) {
	return []string{
		PathClassApp,
		PathClassApplicationPort,
		PathClassConcreteInfra,
		PathClassDomain,
		PathClassDomainErrors,
		PathClassTestMocks,
	}
}
