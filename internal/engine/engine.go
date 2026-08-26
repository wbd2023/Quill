package engine

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/wbd2023/quill/internal/ecosystem/golang"
	"github.com/wbd2023/quill/internal/ecosystem/node"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/external"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/pack/shipped/bindings"
	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/styleguide"
	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

/* ----------------------------------------- Engine Core ---------------------------------------- */

// Engine coordinates repository loading, profile compilation, tool inspection, rule execution,
// installation, coverage, and lock generation for a single repository.
//
// Engine holds only immutable configuration. It does not cache a loaded profile, compiled plan,
// style guide, or toolchain state between operations. Each method loads a fresh snapshot through
// one preparation pipeline.
type Engine struct { // style: allow-package-stutter because: foundational package type
	root           string
	commandRunner  toolchain.CommandRunner
	progressWriter io.Writer
}

// Option configures an Engine.
type Option func(configuration *engineConfiguration) error

type engineConfiguration struct {
	root           string
	commandRunner  toolchain.CommandRunner
	progressWriter io.Writer
}

// New constructs an Engine for the repository rooted at root. The default command runner
// executes local commands.
func New(root string, options ...Option) (engine *Engine, err error) {
	configuration := engineConfiguration{
		root:           root,
		commandRunner:  process.Runner{},
		progressWriter: io.Discard,
	}

	for _, option := range options {
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}

	canonicalRoot, err := workspace.CanonicalRoot(configuration.root)
	if err != nil {
		return nil, err
	}

	return &Engine{
		root:           canonicalRoot,
		commandRunner:  configuration.commandRunner,
		progressWriter: configuration.progressWriter,
	}, nil
}

// WithCommandRunner replaces the command runner used for tool inspection.
func WithCommandRunner(commandRunner toolchain.CommandRunner) (option Option) {
	return func(configuration *engineConfiguration) error {
		configuration.commandRunner = commandRunner
		return nil
	}
}

// WithProgressWriter sets the writer for tool-installation and lock-resolution progress messages.
// The default discards all output.
func WithProgressWriter(writer io.Writer) (option Option) {
	return func(configuration *engineConfiguration) error {
		configuration.progressWriter = writer
		return nil
	}
}

/* ----------------------------------------- Preparation ---------------------------------------- */

// preparedOperation holds the freshly loaded and validated state shared by every repository
// operation. Engine loads it once per operation and never retains it.
type preparedOperation struct {
	config        profile.Profile
	plan          style.Plan
	registry      pack.Registry
	document      styleguide.Document
	externalPacks []pack.Definition
}

// resolvedDrivers holds the check and fix Driver sets for one runner operation.
type resolvedDrivers struct {
	check execution.DriverSet
	fix   execution.DriverSet
}

// prepare loads a fresh repository snapshot for one operation: the Profile,
// the STYLE.md document, the Pack catalog, and the compiled Plan.
// Requirement IDs bound by the Profile are validated against the parsed document
// before the Profile is compiled, so an unknown but syntactically valid
// Requirement ID fails before any Rule or Tool operation.
func (engine *Engine) prepare(
	ctx context.Context,
) (prepared preparedOperation, err error) {
	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	config, err := profile.Load(engine.root)
	if err != nil {
		return preparedOperation{}, err
	}

	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	document, err := styleguide.Load(engine.root, styleguide.Config{
		Filename: config.StyleGuide.Path,
	})
	if err != nil {
		return preparedOperation{}, err
	}

	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	if err = validateRequirementBindings(config, document); err != nil {
		return preparedOperation{}, err
	}

	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	externalPacks, err := external.LoadSources(engine.root, config.PackSources)
	if err != nil {
		return preparedOperation{}, err
	}

	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	registry, err := shipped.ComposeCatalog(externalPacks).Registry(config.EnabledPacks)
	if err != nil {
		return preparedOperation{}, err
	}

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		return preparedOperation{}, err
	}

	if err := ctx.Err(); err != nil {
		return preparedOperation{}, err
	}

	plan, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		return preparedOperation{}, err
	}

	return preparedOperation{
		config:        config,
		plan:          plan,
		registry:      registry,
		document:      document,
		externalPacks: externalPacks,
	}, nil
}

// prepareRun loads the fresh operation snapshot and resolves the executable runner context and
// shipped drivers for a check, fix, inspect, install, or lock operation. Metadata-only operations
// such as Coverage call prepare directly and never construct a runner context, drivers, or
// inspected tools.
func (engine *Engine) prepareRun(
	ctx context.Context,
	scope style.Scope,
) (run execution.RunContext, resolved resolvedDrivers, err error) {
	prepared, err := engine.prepare(ctx)
	if err != nil {
		return execution.RunContext{}, resolvedDrivers{}, err
	}

	config := prepared.config
	if scope == "" {
		scope = config.Repository.DefaultScope
	}

	if !config.Repository.HasScope(scope) {
		return execution.RunContext{}, resolvedDrivers{}, errUnknownScope(scope)
	}

	built := bindings.Build()
	if err = built.Validate(prepared.plan); err != nil {
		return execution.RunContext{}, resolvedDrivers{}, err
	}

	driverSets := resolvedDrivers{
		check: drivers.CheckDrivers(built),
		fix:   drivers.FixDrivers(built),
	}

	layout := workspace.NewLayout(engine.root)
	path := layout.BuildPath(os.Getenv("PATH"), node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := golang.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	runContext := execution.NewRunContext(
		engine.root,
		scope,
		prepared.config,
		prepared.plan,
		prepared.registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
	)

	return runContext, driverSets, nil
}
