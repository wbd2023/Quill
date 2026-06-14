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

## Directory Layout

The tool is a private Go module owned by this repository. The tree is grouped by ownership:
profile model, Pack declaration, Check implementation, execution, output, and test guardrails.

```text
tools/
|-- cmd/
|   `-- style/                         CLI entrypoint for bin/style.
|
`-- internal/
    |-- style/                         Shared rule, scope, diagnostic, and execution vocabulary.
    |-- policy/                        Neutral typed profile policy and repository policy.
    |-- profile/                       Profile facade: load, parse, format, validate, compile.
    |   |-- toml/                      Persisted style.toml schema and policy conversion.
    |   `-- internal/
    |       |-- validation/             Consistency checks for typed profile policy.
    |       |-- effective/              Effective Profile compilation.
    |       `-- profiletest/            Profile-specific test helpers.
    |
    |-- pack/                          Neutral Pack definitions, catalogues, and registries.
    |   `-- shipped/                   Built-in Pack catalogue.
    |       |-- golang/                 Go Pack declaration.
    |       |-- text/                   Text Pack declaration.
    |       |-- bash/                   Bash Pack declaration.
    |       |-- markdown/               Markdown Pack declaration.
    |       |-- project/                Project Pack declaration.
    |       |-- security/               Security Pack declaration.
    |       |-- vocabulary/             Vocabulary Pack declaration.
    |       |-- tool/                  Canonical shipped Tool IDs and capabilities.
    |       `-- bindings/              Shipped Pack Runtime Binding assembly.
    |
    |-- checks/                        Check implementations and Pack-specific policy codecs.
    |   |-- golang/                    Go checker facade and Go check family.
    |   |   |-- analysis/              Shared Go analysis primitives.
    |   |   |-- check/                 Go check IDs.
    |   |   |-- syntax/                AST and syntax checks.
    |   |   |-- structure/             File shape, order, and spacing checks.
    |   |   |-- relationships/         Interface, implementation, and mock checks.
    |   |   |-- architecture/          Go import and layering checks.
    |   |   |-- test/                  Go test hygiene checks.
    |   |   `-- scenarios/             End-to-end Go style scenarios.
    |   |-- text/                     Executable text scanners.
    |   |-- bash/                     Executable Bash scanners.
    |   |-- security/                 Executable security scanners.
    |   |-- vocabulary/               Executable project-term vocabulary scanner.
    |   |-- gopolicy/                 Go Pack Policy codec and validation.
    |   |-- textpolicy/               Text Pack Policy codec and validation.
    |   |-- projectpolicy/            Project Pack Policy codec and validation.
    |   `-- vocabularypolicy/         Vocabulary Pack Policy codec and validation.
    |
    |-- runner/                       Generic rule and fix execution.
    |   `-- drivers/                  Driver facade and execution-family implementations.
    |       |-- command/               File-command execution.
    |       |-- project/               Project-level check execution.
    |       |-- scan/                  Repository scanner execution.
    |       |-- target/                Target command and target check execution.
    |       `-- internal/
    |           |-- runtimebinding/    Runtime Binding contracts and duplicate guards.
    |           `-- commandrun/        Command output handling shared by command drivers.
    |
    |-- toolchain/                    Tool capability, health, command lookup, and versions.
    |-- runtime/                      Command execution and repository-local tool layout.
    |-- installer/                    Pinned Tool installation, downloads, and archives.
    |
    |-- styleguide/                   STYLE.md parsing and hidden metadata extraction.
    |-- requirementid/                STYLE.md requirement ID grammar.
    |-- markers/                      Hidden STYLE.md marker parsing.
    |-- coverage/                     STYLE.md/profile/rule coverage graph.
    |-- report/                       Text and JSON output rendering.
    |   `-- testdata/                 Golden output fixtures.
    |
    |-- cli/                          Command parsing, repo-root detection, and user UX.
    |-- filewalk/                     Repository file collection and generated-file filtering.
    |-- testutil/                     Shared test-only repository helpers.
    |   `-- profiles/                 Repository-shaped profile test setup.
    `-- architecture/                 Test-only architecture boundary guardrails.
```

### Naming Notes

- `internal/policy` is neutral profile policy. It should not know about concrete Packs.
- `checks/<domain>` packages run checks. For example, `checks/text` runs text scanners.
- `checks/<domain>policy` packages decode and validate Pack-specific policy. For example,
  `checks/textpolicy` owns Text Pack Policy, not executable scanners.
- `internal/architecture` is intentionally test-only. It enforces package import boundaries and
  requirement ownership from one place.
- `internal/requirementid` is separate from `styleguide` so profile code can use the requirement
  ID grammar without importing the full STYLE.md parser.
- `internal/profile/internal/profiletest` is profile-specific test support. Shared test helpers
  that are not profile-specific live in `internal/testutil`.
- Report surfaces use `<surface>.go`, `<surface>_text.go`, `<surface>_json.go`, and
  `<surface>_view.go` where those roles exist.
- Multi-role shipped Packs use `execution_ids.go`, `rules.go`, and `file_sets.go`. Small Packs
  stay flat until splitting improves locality.

There is no longer a shell registry, shell-script control plane, nested `tools/style/` module, or
single overloaded package mixing style-platform vocabulary with STYLE.md parsing.
