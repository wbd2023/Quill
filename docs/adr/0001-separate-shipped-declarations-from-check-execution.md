# ADR 0001: Separate shipped declarations from check execution

## Status

Accepted.

This decision originated while Quill was embedded in Ciphera. It has been updated to describe the
standalone implementation and current package names.

## Context

Quill ships reusable Packs while keeping execution infrastructure generic. A Pack must declare its
profile-visible Rules, tools, file-selection defaults, and Pack Policy, but those declarations must
not make generic execution code depend on the Shipped Pack catalogue.

Concrete Checks also have a different responsibility from Pack declarations. Packs define what a
Rule means and which execution identity it selects. Checks observe repository state. Drivers adapt a
resolved Job to a generic execution family. Toolchain code inspects capabilities independently of
which shipped Pack requires them.

Combining these concerns would make Profile vocabulary, Shipped Pack policy, execution mechanics,
and tool health change together. It would also make adding or reorganising a Shipped Pack alter
Driver dependencies.

## Decision

Quill separates these responsibilities:

- `internal/pack/shipped/<pack>` owns Shipped Pack identity and profile-visible Rule declarations.
- `internal/pack/shipped/<pack>/policy` owns that Pack's typed persisted Policy codec.
- `internal/pack/shipped/<pack>/bindings` maps Pack-qualified execution identities to concrete
  scanners, commands, and Checks.
- `internal/pack/shipped/bindings` assembles those Pack-local bindings and Tool-global file
  interpreter mappings into one explicit runtime value.
- `internal/checks` owns concrete repository observations.
- `internal/execution/drivers` owns generic adapters for resolved Job families and the
  `drivers.Bindings` completeness invariant.
- `internal/toolchain` owns generic capability health, inspection, and version detection.
- `internal/installer` owns verified installation of external tools.
- `internal/engine` composes Pack definitions and runtime bindings without moving shipped identities
  into generic execution packages.

`pack/shipped/bindings.Build` is the explicit composition point for shipped runtime behaviour.
Pack-local binding children contain concrete adapters. `drivers.NewBindings` constructs the generic
binding collection, and `Bindings.Validate` confirms every active Rule Job has the required binding.
Drivers do not import shipped Packs, Pack Policy, or concrete Checks.

## Consequences

- Profile language remains independent of Driver implementation details.
- Generic Drivers can be tested without loading Quill's Shipped Pack catalogue.
- Shipped Packs can add Rules and bindings without introducing global registration or `init` side
  effects.
- Check implementations remain reusable across execution families without owning presentation,
  scope, or enforcement policy.
- Tool inspection and installation remain shared platform capabilities rather than per-Pack logic.
- Composition must register every shipped execution identity exactly once. Tests and architecture
  checks must detect missing bindings, duplicate bindings, and forbidden import direction.
- Adding a new execution family may require coordinated changes to the binding facade, shipped
  composition, tests, and architecture documentation.

## Alternatives considered

### Keep shipped declarations and Checks together

Rejected. Pack declarations own profile-visible policy, while Checks own repository observations.
Combining them would couple policy compilation to execution implementation.

### Let generic Drivers import shipped packages

Rejected. Generic execution would silently depend on Quill's default catalogue and could no longer
be reasoned about or tested independently.

### Register shipped behaviour globally

Rejected. Mutable registries and `init` side effects hide composition order, complicate tests, and
allow incomplete runtime state.

### Put generic toolchain behaviour under shipped tools

Rejected. Capability inspection, version detection, and installation are platform concerns shared by
all Packs and operations.

### Call shipped Packs built-ins

Rejected. `shipped` is the established product term and describes distribution ownership. `builtin`
would expose an implementation detail and obscure the distinction between repository policy and
Quill-provided capabilities.
