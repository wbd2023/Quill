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
- Requirement IDs live in hidden `<!-- style: id=... -->` metadata comments instead of the prose
  bullets themselves.
- Hidden `<!-- style: ... -->` metadata comments also declare review-only requirements and other
  machine-only guide metadata.
- Rule packs define checker capabilities; `style.toml` decides which rule packs and rules are
  active.
- Rules map to requirement IDs through `style.toml`, not through implementation code.
- Go style diagnostics use checker-owned diagnostic codes instead of hardcoded STYLE.md IDs.
- Coverage is derived from requirements, not maintained as a hand-written section-status table.

## Implementation

- `cmd/style/`
  - Go CLI entrypoint for `bin/style`.
- `internal/contract/`
  - Shared style-platform contracts, value types, scope/profile enums, executor IDs, and
    Make-surface contracts.
- `internal/rulepack/`
  - Rule-pack registry, builtin rule definitions, tool ownership, scanner IDs, and safe fix specs.
- `internal/executors/`
  - Built-in executor bindings from generic executor/scanner IDs to concrete checks and fixers.
- `internal/styleguide/`
  - STYLE.md parser, requirement model, coverage builder, and exception-marker helpers.
- `internal/profile/`
  - Strict `style.toml` decoding, profile validation, and effective rule compilation.
- `internal/cli/`
  - Command parsing, repository-root resolution, and public CLI UX.
- `internal/runner/`
  - Generic rule/fix execution through injected executors, toolchain inspection, status mapping,
    and profile-owned file-set selection.
- `internal/filewalk/`
  - Repository file collection and generated-file filtering shared by runners and scanners.
- `internal/report/`
  - Text and JSON renderers for checks, coverage, and tool status.
- `internal/rules/go/`
  - Go-specific rule engine and rule entrypoints, with `checks/` for per-file AST passes,
    `order/` for interface and implementation ordering, and `behaviour/` for multi-file
    behaviour tests.
- `internal/rules/repo/`
  - Repository-wide scanners for text, structure, shell, and file-layout rules.
- `internal/runtime/`
  - Installed-tool inspection, repo-local runtime layout, downloads, and subprocess execution.
- `internal/fixtures/`
  - Shared helpers for writing repository-shaped test fixtures.

There is no longer a shell registry, shell-script control plane, nested `tools/style/` module, or
single overloaded package mixing platform contracts with STYLE.md parsing.
