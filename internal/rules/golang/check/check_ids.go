// Package check names the Go checks that a Pack can ask the Go driver to run.
package check

const (
	Comments           = "comments"
	Data               = "data"
	DomainValues       = "domain_values"
	Errors             = "errors"
	FileShape          = "file_shape"
	GuardClauseSpacing = "guard_clause_spacing"
	Logging            = "logging"
	Naming             = "naming"
	Order              = "order"
	Parameters         = "parameters"
	Process            = "process"
	Resources          = "resources"
	Returns            = "returns"
	Security           = "security"
	SwitchCaseSpacing  = "switch_case_spacing"
	Tests              = "tests"
)

func IDs() (ids []string) {
	return []string{
		Comments,
		Data,
		DomainValues,
		Errors,
		FileShape,
		GuardClauseSpacing,
		Logging,
		Naming,
		Order,
		Parameters,
		Process,
		Resources,
		Returns,
		Security,
		SwitchCaseSpacing,
		Tests,
	}
}
