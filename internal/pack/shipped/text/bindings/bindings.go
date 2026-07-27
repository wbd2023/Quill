// Package bindings owns the Text Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Text execution identities (the misspell file interpreter
// and the text repository scanners) to concrete generic drivers. The parent text package stays
// independent of the driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/text"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
)

// Register wires every Text execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerFileInterpreters(bindings)
	registerRepositoryScanners(bindings)
}

func registerFileInterpreters(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Misspell,
		drivers.InterpretPlainText(drivers.ExitFindingsMisspell, "text/spelling/findings"),
	)
}

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(text.ScannerASCII, drivers.CheckASCII())
	bindings.AddRepositoryScanner(
		text.ScannerExceptionMarkers,
		drivers.CheckExceptionMarkers(),
	)
	bindings.AddRepositoryScanner(text.ScannerLineLength, drivers.CheckLineLengths())
	bindings.AddRepositoryScanner(
		text.ScannerMaintenanceMarkers,
		drivers.CheckMaintenanceMarkers(),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaderNames,
		drivers.CheckSectionHeaderNames(text.PackID),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaderDensity,
		drivers.CheckSectionHeaderDensity(text.PackID),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaders,
		drivers.CheckSectionHeaders(text.PackID),
	)
}
