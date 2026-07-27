// Package bindings owns the Markdown Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Markdown execution identities (the markdownlint file
// interpreter) to concrete generic drivers. The parent markdown package stays independent of the
// driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
)

// Register wires every Markdown execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Markdownlint,
		drivers.InterpretPlainText(drivers.ExitFindings, "markdown/markdownlint/findings"),
	)
}
