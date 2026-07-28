# Quill ideal MVP roadmap

[mvp.md](mvp.md) is the canonical ideal first-release target. This document is
the mutable, ordered delivery plan for that target. It records completed work,
remaining work, and dependency order; it does not define a smaller feature cut.

Quill's stable integration surface remains the `quill` command and the
repository-owned `STYLE.md`, `quill.toml`, and `quill.lock` files. All Go
implementation packages remain internal.

## Phase 1: Complete the CLI-first release foundation

Completed:

- Removed the root Go facade and routed `internal/cli` directly to
  `internal/engine`.
- Defined versioned JSON envelopes, stdout purity, structured failures, and
  signal cancellation for machine commands.
- Added tagged release automation that builds, executes, and smoke-tests every
  declared platform archive on its native runner (Linux amd64 and arm64, macOS
  amd64 and arm64, Windows amd64), then verifies archive contents and checksum
  coverage before publishing.

Remaining:

- Add `quill_version` to every machine envelope.
- Publish per-command schema, field, path, position, error, and compatibility
  documentation.
- Keep cancellation, child-process cleanup, stdout purity, exit codes, binary
  archives, checksums, and release smoke tests as release gates.

## Phase 2: Policy integrity and Pack composition

- Validate every Rule `requirement_ids` entry against the loaded `STYLE.md`
  requirement set during Profile preparation.
- Share one validated requirement graph between coverage and execution.
- Decide whether enabling a Pack activates its default Rules while preserving
  consumer-owned enforcement, scope, and requirement bindings.
- Co-locate each shipped Pack's Rule declarations and runtime bindings without
  `init` registration or mutable global registries.
- Give invalid Tool diagnostics one reachable owner rather than competing
  preflight and Driver paths.

## Phase 3: Profile and Tool lifecycle

- Extend repository exclusions to validated repository-relative file and path
  patterns at the file-walk boundary.
- Revisit the flat `[[rules]]` list only when demonstrated authoring failures
  justify a clearer binding and Pack-default model.
- Add version constraints and multiple pinned Tool instances when a real Pack
  needs distinct versions of one Tool.
- Replace the sealed `toolchain.InstallMethod` set with registration only when
  a real fifth installation ecosystem requires it.
- Consolidate shared Pack TOML decoding primitives without creating a
  base-Pack framework.
- Preserve idiomatic vocabulary shorthands such as `ctx`, `err`, `req`, and
  `db`; use generic Driver registries only if real call sites become simpler.

## Phase 4: External Pack data model and runtime

- Define the versioned local Pack manifest and repository-contained source
  layout.
- Validate Pack roots, manifests, IDs, Rule IDs, runtime commands, conflicts,
  and path escapes before execution.
- Compose local external Pack definitions with shipped Pack definitions through
  the existing Profile and engine path.
- Define `quill-pack-v1` JSON Lines requests, diagnostics, completions, and
  protocol failures.
- Run external checks with bounded direct subprocess execution, cancellation,
  output limits, controlled environment inheritance, and captured stderr.
- Migrate built-in and external findings to one documented diagnostic range
  model and render them through the same report path.

## Phase 5: Diagnostic and user-product completion

- Complete the shared diagnostic range model for built-in and external Rules.
- Add deterministic `quill init` presets that never overwrite policy without
  explicit force.
- Add `quill list` and `quill explain` from the loaded compiled snapshot.
- Publish complete CLI, JSON, Profile, Pack-author, and trust-model references.

## Phase 6: Independent adoption proof

- Maintain an independent consumer with a local external Pack, own policy,
  passing and failing fixtures, JSON CI, and a pinned Quill release.
- Cover installation, checks, coverage, JSON automation, protocol failures,
  and Pack security failures end to end.
- Document a reproducible setup that uses the pinned release in consumer CI.

## Phase 7: Public release

- Verify every supported platform and its Tool assumptions before publishing.
- Finalise persisted-format and machine-schema compatibility versions.
- Publish tagged archives, checksums, release notes, installation guidance, and
  persisted-format migration notes.
- Run the complete security, architecture, protocol, consumer, and release
  smoke gates before tagging the ideal MVP.
