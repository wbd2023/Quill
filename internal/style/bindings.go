package style

// RuntimeBindings resolves the Pack-qualified runtime binding identities a compiled Plan
// references. Tool IDs remain global, so file interpreters are keyed by Tool ID: a tool's output
// format is fixed by the tool, not by the Pack that declares a file-command rule.
//
// The completeness validation (profile.ValidateRuntimeBindings) walks an effective Plan and asks
// the RuntimeBindings whether every active execution resolves to exactly one registered binding.
// The generic Driver facade (drivers.Bindings) implements this interface.
type RuntimeBindings interface {
	HasProfileCheck(packID string, check string) (found bool)
	HasRepositoryScanner(packID string, scanner string) (found bool)
	HasTargetCommand(packID string, language string, action string) (found bool)
	HasTargetCheck(packID string, language string, check string) (found bool)
	HasFileInterpreter(toolID string) (found bool)
}
