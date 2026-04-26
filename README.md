# Tooling

This directory is the repo-owned tooling module.

`STYLE.md` remains the human source of truth. `style.toml` is the executable repository profile:
it enables rule packs, binds checker rules to requirement IDs, declares file sets, and supplies
project-specific paths. The Go control plane in `tools/` compiles those two inputs into an
effective rule graph, installs pinned tools from active rule packs, and reports coverage from the
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
- `style.toml` owns active scopes, file sets, language backends, rule bindings, and tool pins.
- Requirement IDs live in hidden `<!-- style: id=... -->` metadata comments instead of the prose
  bullets themselves.
- Hidden `<!-- style: ... -->` metadata comments also declare review-only requirements and other
  machine-only guide metadata.
- Rule packs define checker capabilities; `style.toml` decides which rule packs and rules are
  active.
- Rules map to requirement IDs through `style.toml`, not through implementation code.
- Rule bindings also own project path-class requirements; rule packs do not import project policy.
- Go style diagnostics use checker-owned diagnostic codes instead of hardcoded STYLE.md IDs.
- Coverage is derived from requirements, not maintained as a hand-written section-status table.
- `profile_version = 1` is the first unreleased current schema. There is no legacy schema support.
- Scoped file sets collect from explicit include scopes that overlap the active scope. The profile
  default scope is only the CLI default, not a hidden "widest" scope.

## Boundaries

The dependency direction is:

`contract -> policy -> profile -> cli`

`toolchain -> runtime/cli/report`

`rulepack -> cli/executors`

`styleguide -> coverage -> cli/report`

`runner -> cli/executors`

Production packages must keep these boundaries:

- `contract` imports no internal package.
- `policy` imports only `contract`.
- `profile` imports `contract` and `policy`, not rule packs, runners, executors, rules, or reports.
- `toolchain` imports no project policy, runtime, rules, profile, or reporting package.
- `rulepack` imports `contract` and `toolchain`, not project policy, runtime, rules, or reports.
- `runner` imports no `profile`, `rulepack`, `runtime`, or `report`.
- Concrete rule packages import no `profile`; Go rules import no `rulepack`.
- `report` owns final text and JSON formatting; rules and executors return data.

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
  - Project policy value types and profile vocabulary such as scopes, path classes, language
    backends, naming config, control-plane config, and architecture config.
- `internal/toolchain/`
  - Installed-tool capability and status values, plus status indexing, sorting, and issue helpers.
- `internal/profile/`
  - Strict `style.toml` schema decoding, policy conversion, validation, rendering for fixtures,
    and effective rule compilation.
- `internal/rulepack/`
  - Builtin rule packs, reusable rule/tool capabilities, opaque executor/scanner IDs, Go check
    IDs, and safe fix specs.
- `internal/coverage/`
  - Pure STYLE.md/profile/rule graph coverage assembly.
- `internal/styleguide/`
  - Pure STYLE.md parser, requirement metadata, verification modes, and exception-marker helpers.
- `internal/executors/`
  - Builtin executor bindings from generic executor/scanner IDs to concrete checks and fixers.
- `internal/cli/`
  - Command parsing, repository-root resolution, and public CLI UX.
- `internal/runner/`
  - Generic rule/fix execution through injected executors, status mapping, and policy-owned
    file-set selection.
- `internal/filewalk/`
  - Repository file collection and generated-file filtering shared by runners and scanners.
- `internal/report/`
  - Text and explicit JSON DTO renderers for checks, coverage, and tool status.
- `internal/rules/golang/`
  - Go-specific rule engine and rule entrypoints, with `checks/` for per-file AST passes,
    `order/` for interface and implementation ordering, and `scenarios/` for multi-file
    scenario tests.
- `internal/rules/text/`
  - Text scanners for line length, ASCII, exception markers, maintenance markers, and section
    headers.
- `internal/rules/security/`
  - Security scanners such as committed-secret detection.
- `internal/rules/naming/`
  - Cross-language naming and vocabulary scanners.
- `internal/rules/bash/`
  - Bash-specific script structure, safety, magic-value, and test-hygiene scanners.
- `internal/runtime/`
  - Repo-local layout construction, installed-tool inspection, downloads, installs, and subprocess
    execution through injected toolchain capabilities.
- `internal/fixtures/`
  - Shared helpers for writing repository-shaped test fixtures.

There is no longer a shell registry, shell-script control plane, nested `tools/style/` module, or
single overloaded package mixing platform contracts with STYLE.md parsing.
