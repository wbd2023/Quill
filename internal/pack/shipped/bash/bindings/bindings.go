// Package bindings owns the Bash Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Bash execution identities (scanners and file interpreters)
// to concrete generic drivers. The parent bash package stays independent of the driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/bash"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
)

// Register wires every Bash execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerFileInterpreters(bindings)
	registerRepositoryScanners(bindings)
}

func registerFileInterpreters(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Shellcheck,
		drivers.InterpretPlainText(drivers.ExitFindings, "bash/shellcheck/findings"),
	)
	bindings.AddFileInterpreter(
		tool.Shfmt,
		drivers.InterpretLines(drivers.ExitFindings, "bash/shfmt/findings"),
	)
}

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(
		bash.ScannerMagicValues,
		drivers.CheckBashMagicValues(),
	)
	bindings.AddRepositoryScanner(bash.ScannerSafety, drivers.CheckBashSafety())
	bindings.AddRepositoryScanner(bash.ScannerStructure, drivers.CheckBashStructure())
	bindings.AddRepositoryScanner(
		bash.ScannerTestHygiene,
		drivers.CheckBashTestHygiene(),
	)
}
