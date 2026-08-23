// Package bindings owns the Security Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Security execution identities (the secrets repository
// scanner) to concrete check behaviour. The parent security package stays independent of the
// driver facade and check implementations.
package bindings

import (
	"context"

	checks "github.com/wbd2023/quill/internal/checks/security"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/security"
	"github.com/wbd2023/quill/internal/style"
)

// Register wires every Security execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(security.PackID, security.ScannerSecrets, scanSecrets)
}

// scanSecrets flags committed secrets across the repository scope.
func scanSecrets(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckSecrets(
		context.RepoRoot,
		context.Profile.Repository,
		context.Scope,
	)
}
