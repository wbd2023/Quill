# Tooling

This directory is the repo-owned tooling module.

`STYLE.md` remains the human source of truth. `style.toml` is the executable repository profile:
it enables Packs, binds checker rules to requirement IDs, declares file sets, and supplies
project-specific paths. The Go control plane in `tools/` compiles those two inputs into an
effective rule graph, installs pinned tools from active Packs, and reports coverage from the
same graph.

## Daily Use

Fresh checkout:

- `make style-install`

Primary gate:

- `make lint`

Required-only fast gate:

- `make lint-required`

Safe auto-fixes:

- `make lint-fix`

Maintenance:

- `make style-doctor`
- `make style-coverage`

Structured output:

- `bin/style check --format json`
- `bin/style coverage --format json`
- `bin/style doctor --format json`
- `bin/style help check`

`make style-install` goes through `style install`, so install logic lives in the same Go control
plane as checking, fixing, doctor, and coverage. On a fresh checkout, run it once before the first
`make lint-required` or `make lint`. Installed tools and caches live under the repo-local
`.cache/style/` tree instead of mutating global GOPATH or home-directory tool paths.

`make style` builds `bin/style` directly, and the `make lint*` and `make style-*` targets build it
on demand first, so the style tool is integrated into the repo like the other executables instead
of being launched via ad hoc `go run` plumbing.

## Target Contract

- `make lint`
  - Runs the full strict STYLE.md gate.
- `make lint-required`
  - Runs the required-only STYLE.md gate.
- `make lint-fix`
  - Runs safe STYLE.md auto-fixes where available.
- `make lint-app`
  - Runs the full STYLE.md gate for app scope only.
- `make lint-tools`
  - Runs the full STYLE.md gate for `tools` only.
- `make style-install`
  - Installs or refreshes pinned style tools.
- `make style-doctor`
  - Checks whether pinned style tools are installed and healthy.
- `make style-coverage`
  - Shows STYLE.md automation coverage.

Testing:

- `make test`
  - Runs all tests, including the tooling module.
- `make test-app`
  - Runs application Go tests only.
- `make test-tools`
  - Runs tooling-module tests only.

When `--repo-root` is omitted, the CLI auto-detects the repository root by walking upward until it
finds the configured profile markers, currently `STYLE.md` and `style.toml`.

## Model

- `STYLE.md` is the canonical style guide.
- `style.toml` is the machine-readable project profile.
- `style.toml` owns active scopes, file sets, Targets, rule bindings, and tool pins.
- Requirement IDs live in hidden `<!-- style: id=... -->` metadata comments instead of the prose
  bullets themselves.
- Hidden `<!-- style: ... -->` metadata comments also declare review-only requirements and other
  machine-only guide metadata.
- Packs define checker capabilities; `style.toml` decides which Packs and rules are
  active.
- Rules map to requirement IDs through `style.toml`, not through implementation code.
- Rule bindings also own project path-class requirements; Packs do not import project policy.
- Go style diagnostics use checker-owned diagnostic codes instead of hardcoded STYLE.md IDs.
- Coverage is derived from requirements, not maintained as a hand-written section-status table.
- `profile_version = 1` is the first unreleased current schema. There is no legacy schema support.
- Scoped file sets collect from explicit include scopes that overlap the active scope. The profile
  default scope is only the CLI default, not a hidden "widest" scope.

## Boundaries

The dependency direction is:

`style -> policy -> profile/toml`

`style -> policy -> profile/internal/validation -> profile`

`style -> policy -> pack -> profile/internal/effective -> profile -> cli`

`toolchain -> runtime -> installer -> cli`

`toolchain -> report`

`pack/shipped/<pack> -> pack/shipped -> profile/internal/effective -> profile -> cli`

`pack/shipped/bindings -> runner/drivers -> cli`

`styleguide -> coverage -> cli/report`

`runner -> runner/drivers -> cli`

Production packages must keep these boundaries:

- `style` imports no internal package.
- `policy` imports only `style`.
- `profile/toml` is the persisted `style.toml` codec. It imports profile policy types, not
  loaders, validators, Packs, runners, drivers, checks, or reports.
- `profile/internal/validation` imports profile policy types and style
  vocabulary, not loaders, Packs, runners, drivers, checks, or reports.
- `profile/internal/effective` imports profile policy types, style vocabulary, and neutral Pack
  definitions, not loaders, Shipped Packs, runners, drivers, checks, or reports.
- `profile` is the public facade over profile loading, TOML, validation, Pack default
  resolution, and Effective Profile compilation. It may import neutral Pack registries, but not
  Shipped Packs, runners, drivers, checks, or reports.
- `toolchain` owns Tool capability, health, status, command lookup, and version
  detection. It imports no project policy, runtime, checks, profile, or reporting package.
- `runtime` owns command execution, command environment layout, and no installation
  orchestration or Tool health policy.
- `installer` imports runtime and style tool types, not project policy, checks, profiles,
  reports, or runners.
- `pack` defines neutral Pack definitions, catalogues, and registries.
- `pack/shipped` assembles the Shipped Pack catalogue and may import Shipped Pack modules.
- `pack/shipped/<pack>` modules own declaration-time Pack concepts and may import Check packages,
  Pack-owned policy codecs, and canonical shipped Tool IDs from `pack/shipped/tool`, but not
  runners, drivers, reports, profiles, or installers.
- `pack/shipped/tool` owns reusable shipped tool capabilities, canonical Tool IDs,
  install kinds, and version kinds.
- `pack/shipped/bindings` owns Shipped Pack Runtime Bindings and may import only the
  top-level `runner/drivers` facade, not driver-family subpackages.
- `runner` imports no `profile`, `pack/shipped`, `runtime`, or `report`.
- `runner/drivers` binds generic Execution Kinds to concrete Drivers from explicit
  `drivers.Bindings` without importing Shipped Packs, profiles, reports, or installation packages.
  Its command, project, scan, and target subpackages stay behind the top-level facade.
- Concrete Check packages import no `profile`; Pack Policy packages such as
  `internal/checks/gopolicy`, `internal/checks/textpolicy`, `internal/checks/projectpolicy`,
  and `internal/checks/vocabularypolicy` own typed Pack Policy and avoid Check implementations,
  runners, profiles, and Shipped Packs.
- Go Checks import no `pack/shipped`, and Go Check policy stays separate from Check
  implementations.
- `report` owns final text and JSON formatting; Checks and drivers return data.

## File Shape

The style platform uses balanced granularity:

- Split files when a file owns multiple domain responsibilities.
- Merge tiny glue files when they only contain one helper, alias, or constant and add navigation
  cost.
- Prefer role-named files over broad names such as `types.go`, `helpers.go`, `model.go`, and
  `checks.go`, unless the file is genuinely package-wide.
- Generated and machine-maintained files such as `go.sum` and `package-lock.json` are excluded from
  aesthetic file-shape judgement.

## Implementation

- `cmd/style/`
  - Go CLI entrypoint for `bin/style`.
- `internal/style/`
  - Shared style-platform vocabulary: levels, scopes, check statuses, diagnostics, execution
    results, effective rules, tool policy, and typed execution specs.
- `internal/policy/`
  - Project policy value types and profile vocabulary such as scopes, path roles, Targets,
    Pack configs, and rule bindings.
- `internal/toolchain/`
  - Installed-tool capability and status values, status indexing, command lookup, Tool health
    inspection, version detection, sorting, and issue helpers.
- `internal/runtime/`
  - Command execution, command environment layout, and repository-local tool directories.
- `internal/installer/`
  - Pinned tool installation, downloads, archive extraction, and lockfile validation.
- `internal/profile/`
  - Style Profile facade: loading, parsing, formatting, validation, and Effective Profile
    compilation.
- `internal/profile/toml/`
  - Persisted `style.toml` schema decoding, encoding, and conversion to policy values.
- `internal/profile/internal/validation/`
  - Internal consistency checks for typed Profile policy values.
- `internal/profile/internal/effective/`
  - Effective Profile compilation from resolved Profile policy and Pack registry style definitions.
- `internal/pack/`
  - Neutral Pack definitions, catalogues, registries, and selection validation.
- `internal/pack/shipped/`
  - Shipped Pack catalogue facade and default registry assembly.
- `internal/pack/shipped/<pack>/`
  - Shipped Pack definitions, rule declarations, Tool needs, file-set defaults, and Pack policy
    wiring. Multi-role Packs use `execution_ids.go`, `rules.go`, and `file_sets.go`;
    small Packs stay flat until the split improves locality. Packs reference canonical Tool IDs
    from `internal/pack/shipped/tool/`.
- `internal/pack/shipped/tool/`
  - Reusable shipped Tool IDs, capabilities, install kinds, and version kinds.
- `internal/pack/shipped/bindings/`
  - Shipped Runtime Binding table from scanner IDs, target actions, target-check languages, and
    project check IDs to the generic `runner/drivers` facade.
- `internal/coverage/`
  - Pure STYLE.md/profile/rule graph coverage assembly.
- `internal/styleguide/`
  - Pure STYLE.md parser, requirement metadata, verification modes, and exception-marker helpers.
- `internal/cli/`
  - Command parsing, repository-root resolution, and public CLI UX.
- `internal/runner/`
  - Generic rule/fix execution through injected Drivers, status mapping, and policy-owned file-set
    selection.
- `internal/runner/drivers/`
  - Driver facade that maps Execution Kinds to generic driver families and accepts explicit
    Runtime Bindings for scanner IDs, target actions, target-check languages, and project checks.
- `internal/runner/drivers/internal/runtimebinding/`
  - Runtime Binding function contracts and duplicate-detecting registries shared by the driver
    facade and concrete driver-family packages.
- `internal/runner/drivers/{command,project,scan,target}/`
  - Execution-family Driver implementations owned behind the `runner/drivers` facade.
- `internal/filewalk/`
  - Repository file collection and generated-file filtering shared by runners and scanners.
- `internal/report/`
  - Text and explicit JSON DTO renderers for checks, coverage, and tool status. Report
    surfaces use `<surface>.go`, `<surface>_text.go`, `<surface>_json.go`,
    `<surface>_view.go`, and `<surface>_types.go` where those roles exist.
- `internal/checks/golang/`
  - Go Check facade and package family. The root package walks Go files and reports diagnostics;
    subpackages own check IDs, shared analysis primitives, syntax checks,
    structure checks, relationship checks, architecture checks, test checks, and scenario tests.
- `internal/checks/{gopolicy,textpolicy,projectpolicy,vocabularypolicy}/`
  - Domain-named Pack Policy packages own typed Pack Policy, codecs, and validation for
    Checks that need Profile-supplied policy. These packages use names that match their import
    paths, so callers do not need aliases for ordinary use.
- `internal/checks/text/`
  - Executable Text scanners for line length, ASCII, exception markers, maintenance markers,
    and section headers.
- `internal/checks/security/`
  - Security scanners such as committed-secret detection.
- `internal/checks/vocabulary/`
  - Executable cross-language project-term vocabulary scanner.
- `internal/checks/bash/`
  - Bash-specific script structure, safety, magic-value, and test-hygiene scanners.
- `internal/fixtures/`
  - Shared helpers for writing repository-shaped test fixtures.

There is no longer a shell registry, shell-script control plane, nested `tools/style/` module, or
single overloaded package mixing style-platform vocabulary with STYLE.md parsing.
