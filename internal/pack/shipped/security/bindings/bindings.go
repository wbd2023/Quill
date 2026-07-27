// Package bindings owns the Security Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Security execution identities (the secrets repository
// scanner) to concrete generic drivers. The parent security package stays independent of the
// driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/security"
)

// Register wires every Security execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(security.ScannerSecrets, drivers.CheckSecrets())
}
