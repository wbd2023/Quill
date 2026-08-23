# Quill architecture

Quill is a modular monolith with explicit ownership and one-way dependencies.
Its supported product interfaces are the `quill` command, repository formats,
local External Pack manifest and subprocess protocol, and release artefacts.
All production implementation packages belong under `internal/`.

ADR 0004 is complete: the former root `quill` facade has been removed. No
production Go package exists at the repository root.

## Dependency and runtime flow

Compile-time dependencies flow inward:

```text
cmd/quill
  -> internal/cli
       -> internal/engine
       -> internal/report
       -> internal/style

internal/report
  -> internal/engine, internal/style, internal/coverage, internal/toolchain

internal/engine
  -> capability packages
```

Runtime values return through the presentation and process layers:

```text
engine result -> report DTO or view -> CLI stream and exit policy -> process status
```

A `check` operation follows this order:

1. `workspace` derives repository-local state paths from the discovered repository root.
2. `profile` reads and validates `quill.toml`, then compiles Pack-declared
   Templates and repository policy into an executable `style.Plan`.
3. `engine` selects Rules for the requested scope and mode, inspects their tools,
   explicitly builds and validates runtime bindings, and creates an
   `execution.RunContext` from the immutable operation snapshot.
4. `execution` dispatches each resolved Job to its generic Driver.
5. Drivers invoke external commands or repository Checks and return structured diagnostics.
6. `report` turns the engine result into text or JSON without owning execution policy.

Every engine operation loads a fresh repository snapshot. `Engine` retains only constructor
configuration; it does not cache Profiles, plans, tool status, or operation results.

## Package ownership

- `internal/style` owns shared domain vocabulary: RuleDefinitions, Rules, Templates, Jobs,
  diagnostics, plans, and statuses.
- `internal/profile` owns the Profile model, persisted TOML codec, loading, confinement,
  validation, and compilation into a Plan.
- `internal/pack` owns Pack definitions and resolution. `internal/pack/shipped` owns shipped
  declarations; each Pack's `policy` child owns its typed persisted Policy codec and its `bindings`
  child owns concrete runtime wiring.
- `internal/execution` owns RunContext construction, Rule execution, Driver selection, and file-set
  collection. `internal/execution/drivers` owns generic Job-family dispatch, runtime-binding
  completeness, command execution, and output interpretation.
- `internal/checks` owns repository observations. Checks accept validated Pack Policy values and do
  not own consumer paths, scopes, or enforcement levels.
- `internal/toolchain`, `internal/installer`, and `internal/process` own external-tool discovery,
  verified installation, and bounded process execution.
- `internal/styleguide` and `internal/coverage` own STYLE.md parsing and requirement coverage.
- `internal/report` owns explicit text and JSON protocol views over engine results.
- `internal/workspace`, `internal/filewalk`, and `internal/lockfile` own filesystem layout, bounded
  traversal, and resolved archive state.
- `internal/cli` owns Kong command grammar, repository-root and init-target resolution,
  stdout/stderr discipline, and exit-code mapping. Kong is confined to this package; Quill owns
  its public protocol.
- `internal/engine` is the private application facade and composition coordinator. It owns every
  repository workflow, including safe policy initialization and active-rule explanation.

Architecture tests under `internal/architecture` enforce important import and ownership boundaries.
Update those tests when a deliberate ownership move changes a boundary.

## Consumer boundary

A consuming repository owns:

- `STYLE.md`, including stable requirement IDs;
- `quill.toml`, including scopes, Targets, file sets, policy values, Rule bindings, and tool pins;
- `quill.lock`, including verified per-platform archive hashes;
- Make, CI, or other orchestration that invokes a pinned `quill` command.

Quill ships reusable capabilities, not assumptions about a consumer's directory layout or domain
language. Repository discovery requires both `STYLE.md` and `quill.toml`; automation may pass
`--repository-root` to make that boundary explicit.

## Pack model and selection

A written Requirement needs no Pack. An automatically enforced Requirement needs a Rule binding in
the Profile. Each Rule belongs to exactly one Pack, and Rule IDs are unique across the combined
Shipped and External Pack catalogue. They do not need a Pack-ID prefix.
`enabled` is reserved for the Profile's `[packs]` selection key and cannot be a
Pack ID.

Each Pack owns its declarations: Rules and their Templates, any Quill-managed
Tool capabilities and file-set defaults, and validation of the Pack Policy it
exposes. The catalogue owns global identity validation and selection; it never
interprets Pack-specific policy values.

For every operation, Quill:

1. validates the repository Profile;
2. loads each local External Pack source and validates its manifest and runtime;
3. composes the complete Shipped and External catalogue and validates global
   Pack and Rule identity before selecting anything;
4. selects the Profile's enabled Packs;
5. applies Pack defaults and lets each selected Pack validate its own Policy;
   and
6. compiles active Rule bindings into the executable Plan.

An External Pack receives its raw Policy object in its subprocess request and
validates those Pack-specific values itself. An empty object is valid when the
Pack's policy model permits it. A selected Pack that is unavailable, a Policy
for a disabled Pack, or any catalogue identity conflict fails preparation before
execution.

Prefer an existing truthful Rule over creating another. Repository-specific variation belongs in the
owning Pack's Pack Policy. Add a new reusable capability to a Shipped Pack. Add a capability that is
specific to one repository or organisation as an External Pack. Distinct Rules may overlap
semantically and may bind the same Requirement when they provide different, intentional evidence.

An External Pack is warranted only when the capability and its implementation should remain local to
the repository or organisation. Its manifest and subprocess interface are defined in
[`pack-protocol.md`](pack-protocol.md).

## Change rules

- Add behaviour to the package that owns the concept; `internal/cli` calls the
  private application facade rather than maintaining a second operation path.
- Keep composition explicit. Do not use `init` registration or mutable global registries.
- Validate repository and network input at its boundary, then trust the validated operation model.
- Keep presentation out of Checks and execution policy out of report writers.
- Keep command protocol DTOs explicit, presentation-free, and independent of
  internal implementation types.
