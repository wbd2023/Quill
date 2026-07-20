# ADR 0001: Separate shipped declarations from check execution

## Status

Accepted.

This decision originated while Quill was embedded in Ciphera. It has been updated to describe the
standalone implementation and current package names.

## Context

Quill ships reusable Packs while keeping execution infrastructure generic. A Pack must declare its
profile-visible Rules, tools, file-selection defaults, and Pack Policy, but those declarations must
not make generic execution code depend on the built-in catalogue.

Concrete Checks also have a different responsibility from Pack declarations. Packs define what a
Rule means and which execution identity it selects. Checks observe repository state. Drivers adapt a
resolved Job to a generic execution family. Toolchain code inspects capabilities independently of
which shipped Pack requires them.

Combining these concerns would make Profile vocabulary, built-in policy, execution mechanics, and
tool health change together. It would also make adding or reorganising a shipped Pack alter generic
Driver dependencies.

## Decision

Quill separates these responsibilities:

- `internal/pack/shipped` owns built-in Pack identity and profile-visible Rule declarations.
- `internal/pack/shipped/bindings` maps shipped execution identities to concrete scanners, commands,
  Checks, and file interpreters.
- `internal/checks` owns concrete repository observations and typed Pack Policy codecs.
- `internal/execution/drivers` owns generic adapters for resolved execution families and receives a
  complete `drivers.Bindings` value during composition.
- `internal/pack/shipped/tool` owns the built-in Tool catalogue and capability definitions.
- `internal/toolchain` owns generic capability health, inspection, and version detection.
- `internal/installer` owns verified installation of external tools.
- `internal/engine` composes Pack definitions and runtime bindings without moving shipped identities
  into generic execution packages.

`pack/shipped/bindings.Build` is the explicit composition point for shipped runtime behaviour.
`drivers.NewBindings` constructs the generic binding collection. Drivers do not import shipped Pack
packages or concrete Check packages to discover behaviour implicitly.

## Consequences

- Profile language remains independent of Driver implementation details.
- Generic Drivers can be tested without loading Quill's built-in Pack catalogue.
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
