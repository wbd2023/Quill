package engine

import (
	"context"
	"io"
	"path/filepath"

	"ciphera/tools/internal/ecosystem/golang"
	"ciphera/tools/internal/ecosystem/node"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/pack/shipped/bindings"
	"ciphera/tools/internal/process"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
	"ciphera/tools/internal/workspace"
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

// PackProvider constructs the pack registry and matching drivers after the profile's enabled
// packs have been loaded.
type PackProvider interface {
	Load(
		operationContext context.Context,
		enabledPacks []string,
	) (environment PackEnvironment, loadError error)
}

// PackEnvironment contains pack metadata and the drivers capable of executing the corresponding
// definitions.
type PackEnvironment struct {
	Registry     pack.Registry
	CheckDrivers execution.DriverSet
	FixDrivers   execution.DriverSet
}

// defaultPackProvider wraps shipped packs, shipped bindings, and standard drivers.
type defaultPackProvider struct{}

func (defaultPackProvider) Load(
	_ context.Context,
	enabledPacks []string,
) (environment PackEnvironment, loadError error) {
	registry, err := shipped.DefaultRegistry(enabledPacks)
	if err != nil {
		return PackEnvironment{}, err
	}

	built := bindings.Build()
	return PackEnvironment{
		Registry:     registry,
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

	environment, err := engine.packProvider.Load(operationContext, config.EnabledPacks)
	if err != nil {
		return compiledProfile{}, err
	}

	config, err = pack.ResolvePacks(config, environment.Registry.Packs())
	if err != nil {
		return compiledProfile{}, err
	}

	effective, err := profile.Compile(config, environment.Registry.Definitions())
	if err != nil {
		return compiledProfile{}, err
	}

	return compiledProfile{
		profile:  effective,
		registry: environment.Registry,
	}, nil
}

func (engine *Engine) prepareRunnerContext(
	operationContext context.Context,
	scope style.Scope,
) (context execution.Context, packs PackEnvironment, prepareError error) {
	config, err := profile.Load(engine.repositoryRoot)
	if err != nil {
		return execution.Context{}, PackEnvironment{}, err
	}

	if scope == "" {
		scope = config.Repository.DefaultScope
	}

	if !config.Repository.HasScope(scope) {
		return execution.Context{}, PackEnvironment{}, errUnknownScope(scope)
	}

	packs, err = engine.packProvider.Load(operationContext, config.EnabledPacks)
	if err != nil {
		return execution.Context{}, PackEnvironment{}, err
	}

	config, err = pack.ResolvePacks(config, packs.Registry.Packs())
	if err != nil {
		return execution.Context{}, PackEnvironment{}, err
	}

	compiled, err := profile.Compile(config, packs.Registry.Definitions())
	if err != nil {
		return execution.Context{}, PackEnvironment{}, err
	}

	layout := workspace.NewLayout(engine.repositoryRoot)
	path := layout.BuildPath(node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := golang.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return execution.NewContext(
		engine.repositoryRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		packs.Registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
	), packs, nil
}
