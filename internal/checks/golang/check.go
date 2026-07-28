package golang

// Check selectors.
const (
	CheckComments Check = iota
	CheckData
	CheckDomainValues
	CheckErrors
	CheckFileShape
	CheckGuardClauseSpacing
	CheckLogging
	CheckNaming
	CheckOrder
	CheckParameters
	CheckProcess
	CheckResources
	CheckReturns
	CheckSecurity
	CheckSwitchCaseSpacing
	CheckTests
)

// Check selects one concrete Go repository observation.
type Check uint8
