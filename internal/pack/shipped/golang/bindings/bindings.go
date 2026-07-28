// Package bindings owns the Go Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Go execution identities (the architecture repository
// scanner, the golangci and go_format target commands, and the Go target style checks) to concrete
// check and toolchain behaviour. The parent golang package stays independent of the driver facade
// and check implementations.
package bindings

import (
	checks "github.com/wbd2023/quill/internal/checks/golang"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
)

// Register wires every Go execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerArchitectureScanner(bindings)
	registerTargetCommands(bindings)
	registerTargetChecks(bindings)
}

func registerArchitectureScanner(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(golang.PackID, golang.ScannerArchitecture, scanArchitecture)
}

func registerTargetCommands(bindings *drivers.Bindings) {
	bindings.AddTargetCommand(
		golang.PackID,
		golang.Language,
		golang.TargetActionGolangci,
		lintTargets,
	)
	bindings.AddTargetCommand(
		golang.PackID,
		golang.Language,
		golang.TargetActionGoFormat,
		formatTargets,
	)
}

func registerTargetChecks(bindings *drivers.Bindings) {
	registerTargetCheck(bindings, golang.TargetCheckComments, checks.CheckComments)
	registerTargetCheck(bindings, golang.TargetCheckData, checks.CheckData)
	registerTargetCheck(bindings, golang.TargetCheckDomainValues, checks.CheckDomainValues)
	registerTargetCheck(bindings, golang.TargetCheckErrors, checks.CheckErrors)
	registerTargetCheck(bindings, golang.TargetCheckFileShape, checks.CheckFileShape)
	registerTargetCheck(
		bindings,
		golang.TargetCheckGuardClauseSpacing,
		checks.CheckGuardClauseSpacing,
	)
	registerTargetCheck(bindings, golang.TargetCheckLogging, checks.CheckLogging)
	registerTargetCheck(bindings, golang.TargetCheckNaming, checks.CheckNaming)
	registerTargetCheck(bindings, golang.TargetCheckOrder, checks.CheckOrder)
	registerTargetCheck(bindings, golang.TargetCheckParameters, checks.CheckParameters)
	registerTargetCheck(bindings, golang.TargetCheckProcess, checks.CheckProcess)
	registerTargetCheck(bindings, golang.TargetCheckResources, checks.CheckResources)
	registerTargetCheck(bindings, golang.TargetCheckReturns, checks.CheckReturns)
	registerTargetCheck(bindings, golang.TargetCheckSecurity, checks.CheckSecurity)
	registerTargetCheck(
		bindings,
		golang.TargetCheckSwitchCaseSpacing,
		checks.CheckSwitchCaseSpacing,
	)
	registerTargetCheck(bindings, golang.TargetCheckTests, checks.CheckTests)
}

// registerTargetCheck binds one Go TargetCheck identity to its concrete Check selector. Every
// Pack-owned TargetCheck* identity maps to exactly one selector, keyed by (language, local ID).
func registerTargetCheck(bindings *drivers.Bindings, id string, selector checks.Check) {
	bindings.AddTargetCheck(golang.PackID, golang.Language, id, styleCheck(selector))
}
