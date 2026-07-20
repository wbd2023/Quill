package engine

import (
	"context"
	"io"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/ecosystem/golang"
	"github.com/wbd2023/Quill/internal/ecosystem/node"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers"
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped"
	"github.com/wbd2023/Quill/internal/pack/shipped/bindings"
	"github.com/wbd2023/Quill/internal/process"
	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
	"github.com/wbd2023/Quill/internal/workspace"
)

/* ----------------------------------------- Engine Core ---------------------------------------- */

// Engine coordinates repository loading, profile compilation, tool inspection, rule execution,
// installation, coverage, and lock generation for a single repository.
//
// Engine holds only immutable configuration. It does not cache a loaded profile, compiled plan,
// or toolchain state between operations. Each method loads a fresh snapshot.
type Engine struct {
	repositoryRoot string
	commandRunner  toolchain.CommandRunner
	packProvider   PackProvider
	progressWriter io.Writer
}

// Option configures an Engine.
type Option func(configuration *engineConfiguration) (optionError error)

type engineConfiguration struct {
	repositoryRoot string
	commandRunner  toolchain.CommandRunner
	packProvider   PackProvider
	progressWriter io.Writer
}

// New constructs an Engine for the repository at repositoryRoot. The default command runner
// executes local commands, and the default pack provider uses the shipped packs and shipped
// driver bindings.
func New(repositoryRoot string, options ...Option) (engine *Engine, optionError error) {
	configuration := engineConfiguration{
		repositoryRoot: repositoryRoot,
		commandRunner:  process.Runner{},
		packProvider:   defaultPackProvider{},
		progressWriter: io.Discard,
	}

	for _, option := range options {
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}

	return &Engine{
		repositoryRoot: configuration.repositoryRoot,
		commandRunner:  configuration.commandRunner,
		packProvider:   configuration.packProvider,
		progressWriter: configuration.progressWriter,
	}, nil
}

// WithCommandRunner replaces the command runner used for tool inspection and command execution.
func WithCommandRunner(commandRunner toolchain.CommandRunner) (option Option) {
	return func(configuration *engineConfiguration) (optionError error) {
		configuration.commandRunner = commandRunner
		return nil
	}
}

// WithProgressWriter sets the writer for installation progress messages. The default discards
// all output.
func WithProgressWriter(writer io.Writer) (option Option) {
	return func(configuration *engineConfiguration) (optionError error) {
		configuration.progressWriter = writer
		return nil
	}
}

// WithPackProvider replaces the provider for pack definitions and execution drivers.
func WithPackProvider(packProvider PackProvider) (option Option) {
	return func(configuration *engineConfiguration) (optionError error) {
		configuration.packProvider = packProvider
		return nil
	}
}

/* ---------------------------------------- Pack Provider --------------------------------------- */

// PackProvider supplies Pack definitions separately from their execution Drivers.
type PackProvider interface {
	Definitions(
		operationContext context.Context,
		enabledPacks []string,
	) (definitions PackDefinitions, loadError error)
	Runtime(
		operationContext context.Context,
		enabledPacks []string,
	) (runtime PackRuntime, loadError error)
}

// PackDefinitions contains Pack metadata used to compile a Profile.
type PackDefinitions struct {
	Registry pack.Registry
}

// PackRuntime contains the Drivers used for check and fix operations.
type PackRuntime struct {
	CheckDrivers execution.DriverSet
	FixDrivers   execution.DriverSet
}

// defaultPackProvider supplies shipped Pack definitions and Drivers.
type defaultPackProvider struct{}

func (defaultPackProvider) Definitions(
	_ context.Context,
	enabledPacks []string,
) (definitions PackDefinitions, loadError error) {
	registry, err := shipped.DefaultRegistry(enabledPacks)
	if err != nil {
		return PackDefinitions{}, err
	}

	return PackDefinitions{Registry: registry}, nil
}

func (defaultPackProvider) Runtime(
	_ context.Context,
	_ []string,
) (runtime PackRuntime, loadError error) {
	built := bindings.Build()
	return PackRuntime{
		CheckDrivers: drivers.CheckDrivers(built),
		FixDrivers:   drivers.FixDrivers(built),
	}, nil
}

/* ----------------------------------------- Preparation ---------------------------------------- */

// compiledProfile holds the resolved profile config and compiled execution plan.
type compiledProfile struct {
	profile  profile.EffectiveProfile
	registry pack.Registry
}

func (engine *Engine) loadCompiledProfile(
	operationContext context.Context,
) (compiled compiledProfile, loadError error) {
	config, err := profile.Load(engine.repositoryRoot)
	if err != nil {
		return compiledProfile{}, err
	}

	definitions, err := engine.packProvider.Definitions(operationContext, config.EnabledPacks)
	if err != nil {
		return compiledProfile{}, err
	}

	config, err = pack.ResolvePacks(config, definitions.Registry.Packs())
	if err != nil {
		return compiledProfile{}, err
	}

	effective, err := profile.Compile(config, definitions.Registry.Definitions())
	if err != nil {
		return compiledProfile{}, err
	}

	return compiledProfile{
		profile:  effective,
		registry: definitions.Registry,
	}, nil
}

func (engine *Engine) prepareRunnerContext(
	operationContext context.Context,
	scope style.Scope,
) (context execution.RunContext, runtime PackRuntime, prepareError error) {
	compiled, err := engine.loadCompiledProfile(operationContext)
	if err != nil {
		return execution.RunContext{}, PackRuntime{}, err
	}

	config := compiled.profile.Profile
	if scope == "" {
		scope = config.Repository.DefaultScope
	}

	if !config.Repository.HasScope(scope) {
		return execution.RunContext{}, PackRuntime{}, errUnknownScope(scope)
	}

	runtime, err = engine.packProvider.Runtime(operationContext, config.EnabledPacks)
	if err != nil {
		return execution.RunContext{}, PackRuntime{}, err
	}

	layout := workspace.NewLayout(engine.repositoryRoot)
	path := layout.BuildPath(node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := golang.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return execution.NewRunContext(
		engine.repositoryRoot,
		scope,
		compiled.profile.Profile,
		compiled.profile.Effective,
		compiled.registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
	), runtime, nil
}
