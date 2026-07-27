// Package bindings owns the Project Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Project execution identities (profile checks) to concrete
// generic drivers. The parent project package stays independent of the driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/project"
)

// Register wires every Project execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddProfileCheck(
		project.CheckEnforcementLevels,
		drivers.CheckProfileEnforcementLevels(),
	)
	bindings.AddProfileCheck(
		project.CheckExcludedDirectories,
		drivers.CheckProfileExcludedDirectories(),
	)
	bindings.AddProfileCheck(
		project.CheckCommands,
		drivers.CheckProfileCommands(project.PackID),
	)
}
