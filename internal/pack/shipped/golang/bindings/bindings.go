// Package bindings owns the Go Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Go execution identities (architecture scanner, target
// commands, and the Go target check) to concrete generic drivers. The parent golang package stays
// independent of the driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
)

// Register wires every Go execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerRepositoryScanners(bindings)
	registerTargetBindings(bindings)
}

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(
		golang.ScannerArchitecture,
		drivers.CheckGoArchitecture(golang.PackID),
	)
}

func registerTargetBindings(bindings *drivers.Bindings) {
	bindings.AddTargetCommand(
		golang.TargetActionGolangci,
		drivers.RunGolangci(
			golang.PackID,
			tool.GolangciLint,
			tool.Goimports,
			golang.Language,
		),
	)
	bindings.AddTargetCommand(
		golang.TargetActionGoFormat,
		drivers.RunGoFormat(golang.PackID, tool.Goimports, golang.Language),
	)
	bindings.AddTargetCheck(golang.Language, drivers.CheckGoStyle(golang.PackID, golang.Language))
}
