package golang

const (
	GoCheckComments           = "comments"
	GoCheckData               = "data"
	GoCheckDomainIdentifiers  = "domain_identifiers"
	GoCheckErrors             = "errors"
	GoCheckFileShape          = "file_shape"
	GoCheckGuardClauseSpacing = "guard_clause_spacing"
	GoCheckLogging            = "logging"
	GoCheckNaming             = "naming"
	GoCheckOrder              = "order"
	GoCheckParameters         = "parameters"
	GoCheckProcess            = "process"
	GoCheckResources          = "resources"
	GoCheckReturns            = "returns"
	GoCheckSecurity           = "security"
	GoCheckSwitchCaseSpacing  = "switch_case_spacing"
	GoCheckTests              = "tests"
)

func CheckIDs() (ids []string) {
	return []string{
		GoCheckComments,
		GoCheckData,
		GoCheckDomainIdentifiers,
		GoCheckErrors,
		GoCheckFileShape,
		GoCheckGuardClauseSpacing,
		GoCheckLogging,
		GoCheckNaming,
		GoCheckOrder,
		GoCheckParameters,
		GoCheckProcess,
		GoCheckResources,
		GoCheckReturns,
		GoCheckSecurity,
		GoCheckSwitchCaseSpacing,
		GoCheckTests,
	}
}
