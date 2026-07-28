# Quill ideal MVP

## Status

This document defines Quill's ideal first public release. It is the canonical
product target for the current implementation and release work.

The [roadmap](roadmap.md) sequences its delivery. It is a mutable plan; this
document owns the required product outcome and acceptance criteria.

It is not a greenfield specification. The current Quill implementation is the
baseline. Existing working capabilities are retained unless a later design
decision explicitly replaces them with a better-compatible implementation.

This is deliberately a complete first-release target, not a reduced feature
cut. It includes the current product, every committed roadmap item, and the
external Pack and adoption proof that establish Quill as reusable.

The ideal MVP therefore means:

> A complete, self-hosting public release that preserves Quill's current
> capabilities, proves it is reusable outside its own repository, and supports
> externally authored local Packs without requiring users to modify or rebuild
> Quill.

---

## 1. Product definition

Quill is a reproducible, extensible style-policy platform for software repositories.

It combines:

- a human-authored `STYLE.md`;
- a machine-readable `quill.toml` Profile;
- reusable Packs;
- built-in and external checks;
- external formatters and linters;
- pinned tool installation;
- verified archive hashes in `quill.lock`;
- safe fixes;
- structured diagnostics; and
- requirement coverage reporting.

Quill does not merely run a collection of linters. It compiles repository-owned policy into an
executable Plan and reports both:

1. violations of the policy; and
2. gaps between the written policy and its automated enforcement.

The central product promise is:

> A repository can define what it considers correct, bind those requirements to executable Rules,
> reproduce the required toolchain, run all checks through one command, and see which requirements
> remain automated, review-only, deferred, or uncovered.

---

## 2. MVP outcome

Quill's MVP is complete when an unrelated repository can:

1. install a pinned Quill release;
2. initialise or author `STYLE.md`, `quill.toml`, and `quill.lock`;
3. enable Quill's built-in Packs;
4. install all required pinned tools;
5. run checks and safe fixes;
6. receive stable text and JSON results;
7. inspect requirement coverage;
8. add a local external Pack without changing Quill's source;
9. bind an external Rule to a `STYLE.md` requirement;
10. execute that Rule through a versioned subprocess protocol; and
11. use the same configuration locally and in CI.

The MVP must demonstrate this workflow in at least one repository other than Quill itself.

---

## 3. Existing baseline

The following capabilities already exist and are part of the MVP baseline. They must not be removed
merely to simplify the release.

### 3.1 CLI operations

Quill retains:

```text
quill check
quill fix
quill doctor
quill coverage
quill install
quill lock
quill version
```

These commands remain the primary supported integration surface.

### 3.2 Repository-owned policy

A consuming repository owns:

```text
STYLE.md
quill.toml
quill.lock
```

`STYLE.md` remains the human source of truth.

`quill.toml` remains the executable Profile containing repository policy, enabled Packs, Pack
configuration, scopes, Targets, file sets, tool pins, and Rule bindings.

`quill.lock` remains the verified archive-hash record for archive-installed tools.

### 3.3 Built-in Pack model

Quill retains its built-in Packs, including the current project, text, Markdown, Bash, Go, security,
and vocabulary capabilities.

A Pack may provide:

- Tool capabilities;
- Rule definitions;
- file-set defaults;
- Pack Policy;
- configuration validation;
- check execution; and
- optional fix execution.

The consuming Profile continues to own:

- whether the Pack is enabled;
- whether a Rule is active;
- the Rule's enforcement level;
- the Rule's scope;
- the Rule's requirement bindings; and
- repository-specific Pack Policy.

### 3.4 Compiled execution Plan

Quill continues to load and validate the Profile before compiling Pack declarations and repository
bindings into an executable Plan.

No Rule may execute until the Profile and Plan are valid.

### 3.5 Execution model

Quill retains its current separation between:

- the engine, which coordinates operations;
- Profile compilation;
- Pack declarations;
- execution Jobs;
- Drivers;
- Checks;
- tool inspection and installation;
- bounded process execution; and
- report rendering.

The MVP must preserve the current one-way dependency boundaries and architecture tests.

### 3.6 Tool reproducibility and installation security

Quill retains:

- exact tool pins;
- repository-local installation state;
- version inspection;
- HTTPS-only archive downloads;
- hash verification through `quill.lock`;
- bounded downloads;
- archive traversal protection;
- archive link handling;
- command timeouts;
- bounded process output; and
- post-install verification.

### 3.7 Requirement coverage

Quill retains requirement coverage reporting derived from the written style guide and the compiled
Profile.

Coverage remains a first-class product feature, not an optional reporting add-on.

### 3.8 Safe fixes

Quill retains all currently supported safe fixes.

The MVP does not require a general-purpose source transformation or edit-merging engine. Existing
deterministic fix behaviour remains valid.

### 3.9 Self-hosting

Quill continues to check its own repository with its own `STYLE.md`, `quill.toml`, and `quill.lock`.

Self-hosting is necessary but not sufficient to prove the MVP.

---

## 4. New MVP capabilities

The MVP adds a focused set of capabilities on top of the existing baseline.

The primary addition is a supported external Pack boundary.

Secondary additions improve adoption, discoverability, validation, and release quality.

---

## 5. Local external Packs

### 5.1 Goal

A repository must be able to add a custom Pack without:

- forking Quill;
- editing Quill's source;
- importing Quill's internal Go packages; or
- recompiling the Quill binary.

The first supported external Pack source is a local directory.

Remote registries, dependency resolution, and arbitrary remote executable installation are deferred.

### 5.2 Profile declaration

A consuming Profile may declare one or more local Pack sources:

```toml
[[pack_sources]]
path = ".quill/packs/company"
```

Paths are resolved relative to the repository root.

A source path must:

- remain inside the repository unless explicitly permitted in the future;
- refer to a directory;
- contain a supported Pack manifest;
- use a unique Pack ID; and
- pass validation before any executable is launched.

### 5.3 External Pack layout

A minimal external Pack has the following shape:

```text
.quill/packs/company/
|-- pack.toml
|-- README.md
`-- bin/
    `-- company-quill
```

A Pack may contain fixtures or source code, but Quill only relies on its declared manifest and
runtime executable.

### 5.4 Pack manifest

The initial manifest format is versioned:

```toml
schema_version = 1

[pack]
id = "company"
name = "Company Engineering Policy"
version = "0.1.0"
quill_protocol = "quill-pack-v1"

[runtime]
command = "bin/company-quill"
timeout = "30s"

[[rules]]
id = "company/no-direct-database-access"
name = "No direct database access"
check = "no-direct-database-access"
supports_fix = false
```

The exact persisted field names may change before release, but the manifest must represent:

- manifest schema version;
- Pack identity;
- Pack version;
- required Quill protocol version;
- runtime command;
- execution limits;
- Rule identity;
- runtime check identifier;
- optional fix identifier;
- Tool requirements, where supported;
- file-set requirements, where supported; and
- Pack configuration expectations.

### 5.5 Manifest restrictions

The MVP manifest does not contain arbitrary executable Go values.

It must be decodable and validatable as data.

External Pack manifests must not expose or mirror Quill's entire internal `pack.Definition` type.
The persisted model should contain only concepts required by an external author.

### 5.6 Pack composition

Quill composes:

```text
built-in Pack definitions
+
local external Pack manifests
=
one validated Pack catalogue
```

The existing `PackProvider` seam should be used or adapted to support this composition without
creating a second Profile or execution pipeline.

Built-in and external Packs must participate in the same:

- Pack selection;
- Rule binding;
- scope resolution;
- requirement binding;
- validation;
- execution;
- reporting; and
- coverage calculation.

### 5.7 Conflicts

Profile preparation must reject:

- duplicate Pack IDs;
- duplicate Rule IDs;
- duplicate Tool IDs with incompatible definitions;
- incompatible file-set ownership;
- unsupported manifest schema versions;
- unsupported Pack protocol versions;
- missing runtime executables;
- malformed runtime limits;
- external Rules with no check operation; and
- Pack IDs enabled by the Profile but unavailable from any source.

No external Pack process may execute before these checks pass.

---

## 6. External Pack subprocess protocol

### 6.1 Rationale

Quill's internal Go packages remain private.

The MVP does not publish a Go SDK or use Go's native plugin mechanism.

External Pack checks run as subprocesses using a small, versioned JSON protocol over standard input
and standard output.

This provides:

- language-independent Pack implementations;
- process isolation;
- crash containment;
- timeout enforcement;
- bounded output;
- independent protocol versioning; and
- freedom to continue refactoring Quill's internal packages.

### 6.2 Invocation model

Quill launches the Pack runtime executable for a selected operation.

Quill writes one request object to standard input.

The Pack writes zero or more result records to standard output using JSON Lines.

Diagnostic or protocol records must not be written to standard error. Standard error is reserved for
human-readable runtime debugging and is captured separately.

### 6.3 Check request

A check request contains at least:

```json
{
  "protocol": "quill-pack-v1",
  "operation": "check",
  "repository_root": "/absolute/path/to/repository",
  "pack_id": "company",
  "rule_id": "company/no-direct-database-access",
  "check_id": "no-direct-database-access",
  "scope": "all",
  "files": [
    "internal/service/users.go"
  ],
  "configuration": {}
}
```

Repository files sent to the plugin use repository-relative paths.

The request may later include Target information and execution metadata, but the MVP should avoid
sending Quill's complete internal Plan.

### 6.4 Diagnostic record

A diagnostic record contains at least:

```json
{
  "type": "diagnostic",
  "code": "direct-database-access",
  "message": "Use the repository interface instead of accessing the database directly.",
  "file": "internal/service/users.go",
  "start": {
    "line": 42,
    "column": 5
  },
  "end": {
    "line": 42,
    "column": 18
  },
  "help_url": ""
}
```

The engine supplies or derives:

- Pack ID;
- Rule ID;
- enforcement level;
- selected scope; and
- requirement bindings.

An external process may not override the Profile's enforcement level or claim a different Rule ID.

### 6.5 Completion record

A successful invocation terminates with:

```json
{
  "type": "complete",
  "success": true
}
```

A protocol-level failure may use:

```json
{
  "type": "complete",
  "success": false,
  "error": "configuration field \"allowed_packages\" is required"
}
```

Unexpected process termination, invalid JSON, missing completion, output overflow, or timeout is
reported as a Rule execution error.

### 6.6 Fix operation

External fixes are optional in the MVP.

A Rule declares whether it supports fixing.

Where supported, Quill sends:

```json
{
  "protocol": "quill-pack-v1",
  "operation": "fix",
  "repository_root": "/absolute/path/to/repository",
  "pack_id": "company",
  "rule_id": "company/format-config",
  "fix_id": "format-config",
  "scope": "all",
  "files": [
    "config/example.conf"
  ],
  "configuration": {}
}
```

The MVP may limit external fixes to one of these models:

1. the Pack directly rewrites declared files and reports changed paths; or
2. the Pack returns complete replacement content for individual files.

Quill should choose one explicit model before implementation. It should not introduce a general
overlapping text-edit protocol unless required by a concrete Pack.

### 6.7 Process controls

Quill must apply:

- `context.Context` cancellation;
- a configured timeout;
- a bounded standard-output buffer;
- a bounded standard-error buffer;
- a minimal, deliberate environment;
- repository-root working directory;
- explicit executable resolution; and
- clear error classification.

Quill is not a sandbox. Documentation must state that a local Pack executable runs with the user's
operating-system permissions.

---

## 7. Diagnostic model

### 7.1 Goal

The MVP must support diagnostics produced by both built-in and external Rules through one stable
interchange representation.

### 7.2 Required fields

The internal or reporting-layer diagnostic model should represent:

```go
type Position struct {
    Line   int
    Column int
    Offset int
}

type Diagnostic struct {
    Code    string
    Message string

    File  string
    Start Position
    End   Position

    HelpURL string
}
```

The implementation may keep legacy convenience fields temporarily during migration, but persisted
JSON and plugin protocol formats must have a clear range model.

### 7.3 Position semantics

The protocol and JSON documentation must define:

- whether lines are zero-based or one-based;
- whether columns are zero-based or one-based;
- whether columns count bytes, Unicode code points, or UTF-16 code units;
- whether end positions are inclusive or exclusive;
- how an unknown end position is represented; and
- how diagnostics without a file location are represented.

The recommended MVP contract is:

- one-based lines;
- one-based UTF-8 byte columns or one-based Unicode code-point columns, chosen explicitly;
- zero for unknown values;
- end position exclusive; and
- repository-relative slash-separated paths in persisted output.

### 7.4 Engine-owned metadata

The final rendered finding may include:

- Rule ID;
- Pack ID;
- source Tool or Check;
- enforcement level;
- requirement IDs;
- check status; and
- whether a fix is available.

These values need not all be supplied by the Pack diagnostic itself.

---

## 8. Requirement validation

### 8.1 Existing gap

A syntactically valid requirement ID must not be accepted merely because its text matches the ID
grammar.

### 8.2 Required behaviour

During Profile preparation, Quill must parse the configured `STYLE.md` and validate every active
Rule binding against the discovered requirement set.

The following are errors:

- a Rule references an unknown requirement ID;
- the style guide contains duplicate requirement IDs;
- malformed requirement metadata is present;
- a requirement binding is empty where at least one binding is required by policy; or
- the configured style-guide path cannot be read.

### 8.3 Coverage consistency

Coverage must be calculated from the same validated requirement set used during Profile preparation.

A Rule may never appear to cover a requirement that does not exist.

### 8.4 Error quality

Errors must identify:

- the Rule ID;
- the unknown or duplicate requirement ID;
- the Profile or style-guide source;
- the relevant file path; and
- the corrective action.

---

## 9. `quill init`

### 9.1 Goal

A new repository should be able to adopt Quill without copying Quill's own large self-checking
Profile.

### 9.2 Command

The MVP adds:

```text
quill init
```

Optional presets:

```text
quill init --preset minimal
quill init --preset go
quill init --preset go-bash
```

The exact preset list may be smaller if necessary.

### 9.3 Generated files

`quill init` creates:

```text
STYLE.md
quill.toml
```

It creates `quill.lock` only when the selected preset requires an archive-installed Tool and a valid
lock entry can be generated or resolved.

The command must not overwrite an existing policy file without an explicit force flag.

### 9.4 Preset principles

Generated Profiles must be:

- small;
- readable;
- documented;
- immediately valid;
- suitable for a new repository; and
- significantly simpler than Quill's own self-checking Profile.

A preset should enable only representative Rules.

The generated `STYLE.md` should contain a small set of real requirements with stable IDs.

### 9.5 Non-interactive use

The MVP should support deterministic non-interactive initialisation suitable for tests and scripts.

Interactive prompts are optional.

---

## 10. Discoverability

### 10.1 `quill list`

The MVP adds:

```text
quill list packs
quill list rules
quill list tools
quill list scopes
```

The command reports both available and active entities where appropriate.

`quill list rules` should show at least:

- Rule ID;
- Pack;
- name;
- active or inactive status;
- enforcement level when active;
- scope when active; and
- fix availability.

### 10.2 `quill explain`

The MVP adds:

```text
quill explain <rule-id>
```

A Rule explanation includes:

- Rule ID;
- owning Pack;
- human-readable name;
- check execution category;
- fix execution category or absence;
- required Tools;
- required file set or Target language;
- active binding;
- enforcement level;
- scope;
- requirement IDs;
- relevant Pack configuration; and
- whether the Pack is built-in or external.

Support for explaining requirement IDs is desirable but not required for the first MVP release.

### 10.3 Source of truth

`list` and `explain` must use the same loaded and compiled repository snapshot as other engine
operations. They must not maintain a separate hand-written Rule catalogue.

---

## 11. JSON contract

### 11.1 Versioning

Every JSON command output must contain top-level `schema_version`,
`quill_version`, `command`, and `status` fields. `schema_version` identifies
the immutable machine protocol; `quill_version` identifies the executable that
produced the response.

Example:

```json
{
  "schema_version": 1,
  "quill_version": "v0.1.0",
  "command": "check",
  "status": "ok"
}
```

### 11.2 Stability

Before the first tagged release, Quill must document:

- the JSON schema for each command;
- field meanings;
- path and position semantics;
- status meanings;
- error representation; and
- compatibility expectations during the pre-1.0 period.

### 11.3 Golden tests

Text and JSON output must have golden or equivalent contract tests.

The JSON tests should verify structure and semantics rather than incidental map ordering.

### 11.4 SARIF

SARIF output is a stretch goal, not a required MVP feature.

The diagnostic model should avoid design choices that would make later SARIF output unnecessarily
difficult.

---

## 12. Release and distribution

### 12.1 Supported integration methods

The MVP supports:

- `go install` with an exact semantic version;
- Go 1.24 `tool` directives for Go repositories; and
- downloadable release binaries for the declared platform matrix.

### 12.2 Initial platform matrix

The ideal MVP officially supports:

- Linux amd64;
- Linux arm64;
- macOS amd64;
- macOS arm64; and
- Windows amd64.

WSL2 uses the Linux builds.

The release gate builds and executes every supported platform's archive on its
native runner. It smoke-tests `quill help`, `quill version`, and one repository
check there before publishing, so a platform is not supported merely because
Quill cross-compiles for it.

Executable support and Tool capability are distinct: the `quill` executable is
validated on every supported platform, but a Pack check that requires an
external Tool is only validated where that Tool's installer is supported. The
cross-platform smoke repository check uses a Pack that needs no external Tool,
so it does not claim a Tool works on a platform where the Tool cannot run.

### 12.3 Release artefacts

A release should publish:

- source;
- compressed binaries for supported platforms;
- SHA-256 checksums;
- release notes;
- installation instructions; and
- upgrade or migration notes where persisted formats change.

### 12.4 Version source

The Git tag remains the single Quill version source.

The release process must verify that built binaries report the intended version through:

```text
quill version
```

### 12.5 CI

Release CI should at minimum:

- run the required Quill gate;
- run the complete test suite;
- run `go vet`;
- build each supported binary;
- verify the reported version;
- generate checksums; and
- smoke-test `quill help`, `quill version`, and one repository check.

---

## 13. Independent example consumer

### 13.1 Requirement

The MVP must be demonstrated in a repository other than Quill.

This may be:

- a dedicated example repository;
- Ciphera after adopting the standalone Quill release; or
- a complete repository-shaped integration fixture that is maintained as a real consumer.

A real separate repository is preferred.

### 13.2 Required demonstration

The consumer must have:

```text
STYLE.md
quill.toml
quill.lock
```

It must use:

- different repository paths from Quill;
- at least one built-in Pack;
- at least one external local Pack;
- at least one required Rule;
- at least one recommendation;
- requirement coverage;
- `quill install`;
- `quill check`;
- JSON output in CI; and
- a pinned Quill release.

### 13.3 External Pack example

The example Pack should be small enough to understand but meaningful enough to prove the extension
boundary.

Suitable examples include:

- forbidding direct database access outside repository adapters;
- requiring a project-specific file header;
- validating a repository-specific configuration file;
- enforcing an approved domain vocabulary; or
- checking a custom directory relationship.

The example must include:

- Pack manifest;
- executable implementation;
- tests;
- README;
- one passing fixture; and
- one failing fixture.

---

## 14. Security and trust model

### 14.1 Built-in Tools

Existing archive security properties remain required.

### 14.2 External Packs

External Pack manifests are untrusted repository input.

Their data must be validated before execution.

External Pack executables are trusted code selected by the repository. Quill does not sandbox them.

Quill must:

- show or document the resolved executable path;
- reject unsupported protocol versions;
- apply timeouts;
- bound output;
- avoid silently inheriting unnecessary sensitive environment variables;
- avoid shell interpolation when launching Pack commands;
- use direct process argument construction;
- report crashes and invalid protocol output distinctly; and
- never execute a Pack while merely listing or validating metadata unless execution is required.

### 14.3 Path handling

External Pack paths must be cleaned and validated.

The MVP should reject runtime command paths that escape the Pack root through `..` or symlink
resolution unless a future explicit trust mechanism permits them.

### 14.4 Documentation

The security documentation must clearly distinguish:

- verified downloaded Tool archives;
- local Pack manifests;
- local Pack executables; and
- Quill's lack of a full sandbox.

---

## 15. Performance and determinism

### 15.1 Scope

The MVP does not require a daemon, LSP server, or incremental semantic database.

### 15.2 Required properties

Quill should:

- walk each required file set efficiently;
- avoid unnecessary duplicate Tool execution where compatible jobs can be grouped;
- retain bounded process execution;
- preserve context cancellation;
- produce deterministic Rule and diagnostic ordering;
- keep caches disposable;
- avoid correctness depending on cache presence; and
- avoid retaining stale Profile or Plan state between engine operations.

### 15.3 Concurrency

Parallel Rule execution is optional unless already implemented.

If added, concurrency must be bounded and output ordering must remain deterministic.

---

## 16. Testing requirements

### 16.1 Existing tests

All existing unit, scenario, integration, installer-security, output, and architecture tests remain
part of the release gate.

### 16.2 External Pack tests

The MVP adds tests for:

- valid local Pack loading;
- unknown manifest schema;
- unsupported protocol version;
- duplicate Pack ID;
- duplicate Rule ID;
- missing runtime command;
- path escape attempts;
- successful check execution;
- zero-diagnostic success;
- diagnostic decoding;
- malformed JSON;
- missing completion record;
- plugin-declared failure;
- non-zero process exit;
- timeout;
- cancellation;
- output truncation;
- standard-error capture; and
- external Rule coverage.

### 16.3 Initialisation tests

`quill init` tests must cover:

- each supported preset;
- valid generated files;
- refusal to overwrite;
- deterministic output;
- explicit force behaviour, if supported; and
- successful `quill check` after initialisation and installation.

### 16.4 Requirement validation tests

Tests must cover:

- known requirement binding;
- unknown but syntactically valid ID;
- duplicate style-guide ID;
- malformed metadata;
- missing style-guide file; and
- coverage consistency.

### 16.5 Release smoke tests

At least one test must build the real CLI and run it against the independent example consumer.

---

## 17. Documentation requirements

The MVP must include:

### 17.1 User documentation

- installation;
- repository initialisation;
- Profile overview;
- running checks;
- applying safe fixes;
- tool installation and locking;
- doctor output;
- requirement coverage;
- JSON usage;
- CI integration; and
- supported platforms.

### 17.2 Pack author documentation

- external Pack directory layout;
- manifest schema;
- protocol versioning;
- check request;
- fix request, if supported;
- diagnostic schema;
- path and position semantics;
- process lifetime;
- cancellation;
- timeouts;
- output limits;
- configuration;
- local testing; and
- security implications.

### 17.3 Rule reference

Built-in Rules should have discoverable documentation containing:

- Rule ID;
- purpose;
- Pack;
- expected configuration;
- check behaviour;
- fix behaviour;
- examples; and
- limitations.

The initial release may generate part of this reference from Pack declarations.

### 17.4 Architecture documentation

The architecture guide must explain:

- the engine facade;
- fresh-snapshot operation loading;
- Profile compilation;
- Pack declarations;
- local Pack composition;
- execution Jobs and Drivers;
- subprocess Pack execution;
- reporting;
- installation and locking; and
- public versus internal compatibility boundaries.

---

## 18. Explicit non-goals

The following are not required for this MVP.

### 18.1 No central Pack registry

External Packs are local.

Git-sourced Packs may be considered later, but they must not block the MVP.

### 18.2 No Pack dependency resolver

External Packs do not declare arbitrary transitive Pack dependencies.

### 18.3 No Pack ecosystem lockfile yet

`quill.lock` continues to lock verified archive hashes for Tools.

It does not yet need to resolve external Pack versions or dependency graphs.

### 18.4 No public Go SDK

All Go implementation packages remain internal.

The CLI, repository files, Pack manifest, and subprocess protocol are the supported integration
surfaces.

### 18.5 No native dynamic plugins

Quill does not use Go plugins, shared libraries, or unstable native ABIs.

### 18.6 No WebAssembly runtime

A future version may add permissioned WebAssembly Packs, but subprocess Packs are sufficient.

### 18.7 No remote Pack execution

Quill does not download and execute arbitrary remote Pack binaries in the MVP.

### 18.8 No LSP or editor daemon

Editor integration and incremental checking are later work.

### 18.9 No declarative syntax-query language

Tree-sitter or another declarative structural Rule format may be added later.

### 18.10 No universal semantic analysis model

Language-specific built-in Checks may use official compiler or parser APIs.

Quill does not create a language-neutral semantic database in the MVP.

### 18.11 No general fix-composition engine

The MVP preserves built-in safe fixes and supports only a deliberately narrow external fix contract.

### 18.12 No distributed or cloud execution

All work runs on the local machine or CI worker.

---

## 19. Acceptance criteria

Quill's MVP is complete when all criteria below are satisfied.

### 19.1 Existing capability preservation

- All current public commands remain available.
- Existing Quill self-checking behaviour remains operational.
- Existing built-in Packs and representative Rules remain operational.
- Existing safe fixes remain operational.
- Tool installation remains pinned and repository-local.
- Archive tools remain verified through `quill.lock`.
- Requirement coverage remains operational.
- Text and JSON reporting remain operational.
- Architecture tests remain operational.

### 19.2 External Pack support

- A Profile can declare a local Pack source.
- Quill validates a versioned Pack manifest.
- The external Pack appears in normal Pack resolution.
- An external Rule can be enabled and bound through ordinary `[[rules]]` configuration.
- An external Rule can bind to a `STYLE.md` requirement.
- Quill invokes the Pack through the versioned subprocess protocol.
- The Pack can return structured diagnostics.
- Diagnostics are rendered with built-in diagnostics.
- A failed or timed-out Pack does not crash Quill.
- The external Rule appears correctly in coverage.
- No Quill source modification or rebuild is required.

### 19.3 Policy integrity

- Every active Rule requirement binding is validated against `STYLE.md`.
- Duplicate requirement IDs fail.
- Unknown requirement IDs fail before execution.
- Coverage uses the same validated requirement graph.
- Errors identify the responsible Rule and requirement.

### 19.4 Onboarding

- `quill init` creates a valid minimal repository policy.
- At least one preset can pass `quill check` after required Tools are installed.
- Existing files are not overwritten without explicit consent.
- Generated policy files are concise and documented.

### 19.5 Discoverability

- Users can list available Packs.
- Users can list available and active Rules.
- Users can list pinned Tools and configured scopes.
- Users can explain an active Rule.
- External Packs and Rules are identified as external.

### 19.6 Automation contract

- JSON output is schema-versioned.
- JSON schema semantics are documented.
- JSON output has contract tests.
- Exit-code behaviour is documented and tested.

- Every envelope includes the producing Quill version.

### 19.7 Distribution

- A semantic-versioned release can be installed through `go install`.
- Linux amd64 and arm64, macOS amd64 and arm64, and Windows amd64 binaries are
  published.
- Release artefacts have SHA-256 checksums.
- `quill version` and every JSON envelope report the release version.
- Release CI builds, verifies, checksums, and smoke-tests every supported
  artefact.

### 19.8 Independent adoption

- At least one repository other than Quill uses the release.
- That repository uses its own Profile and style guide.
- It uses at least one built-in Pack.
- It uses at least one external local Pack.
- It runs Quill in CI through a pinned version.
- Its setup is documented and reproducible.

### 19.9 Lifecycle and maintainability

- Profiles support version constraints and multiple pinned instances when a
  Pack requires distinct versions of the same Tool.
- Tool installation remains typed until a real fifth installation ecosystem
  justifies explicit installer registration.
- Pack activation, default Rule selection, Rule declarations, and runtime
  bindings have one explicit, tested ownership model.
- Toolchain preflight and execution Drivers have one owner for invalid-tool
  diagnostics; no diagnostic path is unreachable.
- Repository exclusions support the validated repository-relative path and file
  patterns required by real Profiles, with filtering kept at the file-walk
  boundary.
- The persisted `[[rules]]` shape is reconsidered when demonstrated authoring
  failures justify a clearer model for bindings and Pack defaults.
- Shared TOML decoding primitives are consolidated without creating a
  base-Pack framework.
- The shipped vocabulary Pack preserves established language idioms.
- Driver binding registries remain typed unless a generic registry makes real
  call sites simpler.

---

## 20. Delivery roadmap

[roadmap.md](roadmap.md) is the single mutable, ordered delivery plan for this
MVP. It records completed work, remaining phases, and their dependencies. A
roadmap update must preserve this document's outcome, non-goals, and acceptance
criteria unless a deliberate product decision changes them.

## 21. Stretch goals

The following may be included if they do not delay the core acceptance criteria:

- SARIF output;
- `quill explain requirement:<id>`;
- Git-sourced Pack directories pinned by commit;
- external Pack checksum entries;
- a dry-run external fix mode;
- Pack manifest JSON Schema;
- generated built-in Rule reference;
- shell completion;
- more `quill init` presets; and
- parallel Rule execution with deterministic output.

Stretch goals do not redefine MVP completion.

---

## 22. Final product statement

The Quill MVP is:

> A self-hosting, reproducible, and externally extensible style-policy platform that compiles
> repository-owned written standards into executable checks, combines built-in semantic Rules with
> pinned external Tools, installs those Tools securely, applies safe fixes, validates requirement
> traceability, reports automation coverage, and supports locally authored third-party Packs through
> a stable file-and-subprocess boundary.

The decisive proof of the MVP is not the number of built-in Rules.

It is that another repository can adopt a pinned Quill release, define its own policy, add a custom
Pack without rebuilding Quill, and receive reproducible checks and coverage through the same
supported workflow used by Quill itself.
