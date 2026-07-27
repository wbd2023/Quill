# Quill architecture

Quill is a modular monolith with explicit ownership and one-way dependencies.
Its supported product boundaries are the `quill` command and repository file
formats. All production implementation packages belong under `internal/`.

ADR 0004 is complete: the former root `quill` facade has been removed. No
production Go package exists at the repository root.

## Runtime flow

```text
cmd/quill
  -> internal/cli
  -> internal/engine
       -> workspace and repository loading
       -> profile validation and compilation
       -> shipped Pack definitions and bindings
       -> toolchain inspection and installation
       -> execution Drivers and Checks
  -> internal/report
```

A `check` operation follows this order:

1. `workspace` derives repository-local state paths from the discovered repository root.
2. `profile` loads `quill.toml`, validates consumer policy, and resolves enabled Pack definitions
   into an executable `style.Plan`.
3. `engine` selects Rules for the requested scope and mode, inspects their tools, and creates an
   `execution.RunContext` from the immutable operation snapshot.
4. `execution` dispatches each resolved Job to a Driver supplied by shipped runtime bindings.
5. Drivers invoke external commands or repository Checks and return structured diagnostics.
6. `report` turns the engine result into text or JSON without owning execution policy.

Every engine operation loads a fresh repository snapshot. `Engine` retains only constructor
configuration; it does not cache Profiles, plans, tool status, or operation results.

## Package ownership

- `internal/style` owns shared domain vocabulary: Rules, Jobs, diagnostics, plans, and statuses.
- `internal/policy` owns the decoded consumer Profile model.
- `internal/profile` owns Profile loading, validation, Pack-policy resolution, and compilation.
- `internal/pack` owns Pack definitions; `internal/pack/shipped` composes shipped Pack declarations
  and delegates Pack-local runtime registration to explicit child binding packages.
- `internal/execution` owns RunContext construction, Rule execution, Driver selection, and file-set
  collection.
- `internal/execution/drivers` adapts resolved Jobs to commands and Checks.
- `internal/checks` owns repository observations and Pack-specific policy codecs. Checks do not own
  consumer paths, scopes, or enforcement levels.
- `internal/toolchain`, `internal/installer`, and `internal/process` own external-tool discovery,
  verified installation, and bounded process execution.
- `internal/styleguide`, `internal/coverage`, and `internal/report` own STYLE.md parsing,
  requirement coverage, and CLI presentation respectively.
- `internal/workspace`, `internal/filewalk`, and `internal/lockfile` own filesystem layout, bounded
  traversal, and resolved archive state.
- `internal/cli` owns argument parsing, stdout/stderr discipline, and exit-code mapping.
- `internal/engine` is the private application facade and composition coordinator.

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
`--repo-root` to make that boundary explicit.

## Change rules

- Add behaviour to the package that owns the concept; `internal/cli` calls the
  private application facade rather than maintaining a second operation path.
- Keep composition explicit. Do not use `init` registration or mutable global registries.
- Validate repository and network input at its boundary, then trust the validated operation model.
- Keep presentation out of Checks and execution policy out of report writers.
- Keep command protocol DTOs explicit, presentation-free, and independent of
  internal implementation types.
