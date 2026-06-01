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

`contract -> policy -> profile/{toml,validation,effective} -> profile -> cli`

`toolchain -> runtime -> installer -> cli`

`toolchain -> report`

`pack/builtin -> runner/drivers -> cli`

`styleguide -> coverage -> cli/report`

`runner -> runner/drivers -> cli`

Production packages must keep these boundaries:

- `contract` imports no internal package.
- `policy` imports only `contract`.
- `profile/toml` imports profile policy types, not loaders, Packs, runners, drivers, rules,
  or reports.
- `profile/validation` imports profile policy types and contracts, not loaders, Packs,
  runners, drivers, rules, or reports.
- `profile/effective` imports profile policy types and contracts, not loaders, Packs, runners,
  drivers, rules, or reports.
- `profile` is a facade over profile loading, TOML, validation, and effective compilation. It does
  not import Packs, runners, drivers, rules, or reports.
- `toolchain` imports no project policy, runtime, rules, profile, or reporting package.
- `runtime` owns command execution, tool inspection, environment layout, and no installation
  orchestration.
- `installer` imports runtime and tool contracts, not project policy, rules, profiles, reports, or
  runners.
- `pack` defines neutral Pack definitions, catalogues, and registries.
- `pack/builtin` assembles built-in Packs and may import rule packages and Pack-owned policy codecs.
- `runner` imports no `profile`, `pack/builtin`, `runtime`, or `report`.
- `runner/drivers` binds generic executor IDs to concrete checks and commands without importing
  profile, report, or installation packages.
- Concrete rule packages import no `profile`; Go rules import no `pack/builtin`, and Go
  rule policy stays separate from rule implementations.
- `report` owns final text and JSON formatting; rules and drivers return data.

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
- `internal/contract/`
  - Shared style-platform contracts: levels, scopes, check statuses, diagnostics, execution
    results, effective rules, tool policy, and typed execution specs.
- `internal/policy/`
  - Project policy value types and profile vocabulary such as scopes, path roles, Targets,
    Pack configs, and rule bindings.
- `internal/toolchain/`
  - Installed-tool capability and status values, plus status indexing, sorting, and issue helpers.
- `internal/runtime/`
  - Command execution, tool inspection, environment layout, and version detection.
- `internal/installer/`
  - Pinned tool installation, downloads, archive extraction, and lockfile validation.
- `internal/profile/`
  - Style profile facade: loading, parsing, formatting, validation, and effective rule compilation.
- `internal/profile/toml/`
  - Persisted `style.toml` schema decoding, encoding, and conversion to policy values.
- `internal/profile/validation/`
  - Internal consistency checks for typed profile policy values.
- `internal/profile/effective/`
  - Compilation from typed profile policy and rule definitions to effective runtime contracts.
- `internal/pack/`
  - Neutral Pack definitions, catalogues, registries, and selection validation.
- `internal/pack/builtin/`
  - Built-in Pack catalogue, reusable rule/tool capabilities, opaque executor/scanner IDs,
    pack-owned policy defaults, and safe fix specs.
- `internal/coverage/`
  - Pure STYLE.md/profile/rule graph coverage assembly.
- `internal/styleguide/`
  - Pure STYLE.md parser, requirement metadata, verification modes, and exception-marker helpers.
- `internal/cli/`
  - Command parsing, repository-root resolution, and public CLI UX.
- `internal/runner/`
  - Generic rule/fix execution through injected executor functions, status mapping, and policy-owned
    file-set selection.
- `internal/runner/drivers/`
  - Built-in Drivers that map generic executor/scanner IDs to concrete checks and fixers.
- `internal/filewalk/`
  - Repository file collection and generated-file filtering shared by runners and scanners.
- `internal/report/`
  - Text and explicit JSON DTO renderers for checks, coverage, and tool status.
- `internal/rules/golang/`
  - Go rule facade and package family. The root package walks Go files and reports diagnostics;
    subpackages own check IDs, Go pack policy, shared analysis primitives, syntax checks,
    structure checks, relationship checks, architecture checks, test checks, and scenario tests.
- `internal/rules/text/`
  - Text scanners for line length, ASCII, exception markers, maintenance markers, and section
    headers.
- `internal/rules/security/`
  - Security scanners such as committed-secret detection.
- `internal/rules/vocabulary/`
  - Cross-language project-term vocabulary scanners.
- `internal/rules/bash/`
  - Bash-specific script structure, safety, magic-value, and test-hygiene scanners.
- `internal/fixtures/`
  - Shared helpers for writing repository-shaped test fixtures.

There is no longer a shell registry, shell-script control plane, nested `tools/style/` module, or
single overloaded package mixing platform contracts with STYLE.md parsing.
