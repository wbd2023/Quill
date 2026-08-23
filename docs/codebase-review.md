# Quill whole-codebase review roadmap

## Purpose and authority

This document owns Quill's systematic whole-codebase review programme: its
coverage, dependency order, evidence standard, phase status, and re-review
triggers. It is for maintainers reviewing the repository by hand or with
review assistance.

This is a design-first review programme. It begins with Quill's product model,
domain language, module ownership, interfaces, and core composition. It then
moves outward to package implementation quality, operational boundaries,
public integration, and delivery assurance.

It is not a product specification, implementation backlog, architecture
reference, protocol reference, security policy, ADR, or append-only audit log.
Those concerns retain their canonical documents. If review exposes an error in
one of them, correct that owner and record only the finding disposition here.

The sequence is intentional. A design review should establish whether Quill is
organised around the right enduring concepts before spending effort on details
that could disappear in a better termination state. A trust-boundary review is
still required, but it comes after the architecture and its responsibilities
are understood.

## Canonical sources

- **Product and distribution:** [product.md](product.md). Confirm what Quill
  promises to consumers.
- **First-release scope:** [mvp.md](mvp.md). Confirm intended scope; do not
  treat it as an audit log.
- **MVP delivery plan:** [roadmap.md](roadmap.md). Route product-scope and
  delivery dependencies here.
- **Package ownership and runtime flow:** [architecture.md](architecture.md).
  Check implementation ownership and import direction.
- **Machine CLI protocol:** [cli-protocol.md](cli-protocol.md). Check JSON,
  streams, errors, exits, and cancellation.
- **Durable architecture decisions:** [adr/](adr/). Check accepted constraints
  and record superseding decisions there.
- **Trust model and vulnerability reporting:** [../SECURITY.md](../SECURITY.md).
  Check trust boundaries and security regression expectations.
- **Repository policy:** [../STYLE.md](../STYLE.md) and
  [../quill.toml](../quill.toml). Check self-hosted policy and its executable
  interpretation.
- **Contributor workflow:** [../CONTRIBUTING.md](../CONTRIBUTING.md). Check
  expected verification and documentation maintenance.
- **Small tactical work:** [../TODO.md](../TODO.md). Route bounded,
  repository-local work here.

## Baseline, scope, and fixed exclusions

The initial review order was established on 2026-07-31. A reviewer beginning a
phase must record the source revision or other reproducible baseline, the
specific paths and symbols traced, and the observed verification evidence.
Existing worktree changes may be unrelated. Review must not reset, clean,
discard, or overwrite them merely to establish a baseline.

### In scope

- `cmd/quill`, every production `internal/` package, and their consumers.
- Root policy and persisted formats: `STYLE.md`, `quill.toml`, and `quill.lock`.
- Tests, test fixtures, architecture tests, release workflow, CI, Make targets,
  module metadata, and durable documentation.
- Success, invalid-input, operational-failure, cancellation, and
  platform-specific paths for every supported operation.

### Pinned CLI-boundary exclusion

ADR 0004 and ADR 0005 are complete and enforced. The following are out of
scope unless an explicit new requirement, verified defect or security
regression, or superseding ADR reopens them:

- restoring a root `package quill` facade or any public Go library;
- adding a second orchestration path, compatibility alias, or generic command
  dispatcher;
- replacing Kong, moving it outside `internal/cli`, or redesigning the current
  CLI package boundary; and
- making `internal/engine` depend on `internal/cli` or `internal/report`.

`cmd/quill` and `internal/cli` remain in scope for caller tracing, protocol,
security, and regression review. The exclusion prevents architecture churn; it
does not suppress evidence of a real behavioural defect.

## Review method and evidence standard

### Design questions

Apply these questions before proposing any structural change:

1. **Responsibility:** What durable concept does this module own? Is that
   concept named in Quill's vocabulary, or is the module merely preserving an
   implementation accident?
2. **Interface:** What must callers know - types, invariants, ordering,
   failures, configuration, and performance - to use the module correctly? Is
   that interface smaller than the behaviour it exposes?
3. **Depth:** Does the module hide meaningful complexity and provide leverage to
   several callers, or is it a thin pass-through that spreads knowledge?
4. **Seam:** Is the boundary where behaviour actually varies? One concrete
   implementation alone does not justify a public-looking interface, registry,
   or framework.
5. **Ownership:** Is each policy decided in one package? Are imports and data
   flow aligned with that ownership?
6. **Termination state:** Ignoring migration effort, would this still be the
   cleanest shape for Quill's actual product? If not, name the coherent final
   shape before suggesting incremental edits.
7. **Locality:** Can a future maintainer understand, change, and verify the
   behaviour within one module and its immediate consumers, or must they trace
   incidental plumbing across the repository?

### Evidence standard

A phase review must trace both directions across its seam:

1. Identify producers, consumers, concrete values, errors, and cancellation
   signals crossing the seam.
2. Read the package implementation, its tests, callers, consumers, related
   configuration, and the document that owns its public or architectural
   contract.
3. Inspect the normal path plus the failure and cancellation paths that matter
   to the module's interface.
4. At filesystem, network, archive, subprocess, and persisted-data boundaries,
   identify untrusted input, the validation point, the trusted model created
   afterward, and the consequence if validation fails.
5. Inspect tests as executable contracts. Prefer behaviour, invariants, and
   error cases over implementation-shaped assertions.
6. Separate observed facts, design judgements, inferences, and unresolved
   questions. A finding must cite files, symbols, a reproducible scenario or
   test, and an owner for its disposition.

Each phase record should contain:

- status: not started, in progress, blocked, complete, or reopened;
- baseline: revision or reproducible source state reviewed;
- scope: packages, formats, commands, and platform paths covered;
- interface: responsibilities, callers, consumers, invariants, and error modes;
- evidence: tests, scenarios, commands, release artefacts, or observations;
- findings: short ID, severity/risk, disposition, and canonical owner; and
- re-review trigger: the material change that invalidates the evidence.

Do not paste prompts, raw tool transcripts, generated reports, or an unbounded
finding history into this document. Keep the latest useful status record and
route durable outcomes to their correct owner.

## System model to review

Quill is a CLI-first modular monolith. Its supported product interfaces are the
`quill` executable, repository formats, versioned JSON, the local External
Pack manifest and subprocess protocol, and release artefacts, not an
in-process public Go API.

The central composition flow is:

```text
cmd/quill
  -> internal/cli
       -> internal/engine
            -> profile + styleguide + Pack composition
            -> execution + Drivers + Checks
            -> toolchain + installer + process + filesystem
       -> internal/report
```

Every engine operation begins with a fresh repository snapshot:

```text
repository root
  -> Profile load and validation
  -> STYLE.md load and requirement binding validation
  -> external Pack source loading
  -> shipped and external Pack catalogue composition
  -> Pack policy resolution and Profile Plan compilation
  -> runtime binding validation
  -> operation-specific execution or metadata work
```

The review must trace these representative flows without copying their
canonical architecture description:

1. **Check and fix:** selection -> tool inspection -> run context -> Job and
   Driver dispatch -> diagnostic/status result -> report -> public exit policy.
2. **Tool lifecycle:** capability resolution -> doctor/install/lock -> lockfile
   hash lookup -> download -> checksum -> archive extraction -> replacement ->
   inspection.
3. **Metadata:** coverage, list, and explain use the same fresh compiled
   snapshot without launching tools unnecessarily.
4. **Safe bootstrap:** init resolves a target, refuses clobbering policy files,
   writes deterministically, rolls back partial state, and honours cancellation.
5. **Product shell:** signals, grammar, help, streams, JSON/text errors, exits,
   build metadata, CI, archive production, checksums, and publication.

## Ordered review phases

### 0. Product model, domain language, and accepted constraints

**Status:** Complete.

**Baseline and evidence (2026-07-31):** reviewed the current worktree without
assuming it is clean. Traced the public-contract documents, `CONTEXT.md`,
`STYLE.md`, root `quill.toml`, ADRs 0001, 0004, 0005, and 0006, the
`internal/style`, `internal/profile`, `internal/pack`, external-Pack,
`internal/engine`, CLI, report, and architecture-guard seams. ADR 0006
resolves the local External Pack public-boundary question. The Profile,
Pack, Rule, Tool, Target, Driver, Check, and Diagnostic model is coherent;
remaining terminology corrections belong to their canonical owners and do not
reopen the pinned CLI boundary.

**Re-review trigger:** material change to the product documents, `CONTEXT.md`,
accepted ADRs, or a public CLI, repository-format, or External Pack contract.

**Why first:** no design judgement is useful until the reviewer understands what
Quill is for, what it deliberately is not, and the vocabulary that should shape
its code.

**Review paths:**

- `README.md`, `docs/product.md`, `docs/mvp.md`, `docs/architecture.md`,
  `docs/cli-protocol.md`, `SECURITY.md`, and `docs/adr/`.
- `CONTEXT.md`, `STYLE.md`, and root `quill.toml` as the repository's stated
  language and executable policy.
- `internal/style`, `internal/profile`, and the package inventory.

**Review questions:**

- Is the CLI-first, language-neutral product boundary still the cleanest model?
- Are the core terms consistent across docs, Profile, types, command output,
  Pack declarations, and tests?
- Do Rule, Pack, Profile, Tool, Target, Driver, Check, scope, requirement, and
  diagnostic each name one distinct concept?
- Are public protocol and persisted-format compatibility obligations separated
  from private Go implementation details?
- Do ADR 0001, ADR 0004, ADR 0005, and ADR 0006 still describe the constraints
  that matter, without preserving obsolete shapes through compatibility residue?

**Exit criteria:** explain Quill's product model, its public interfaces, its
non-goals, its key terms, and the fixed CLI constraints without referring to
implementation trivia.

### 1. Architectural shape, ownership, and package topology

**Status:** Complete.

**Baseline and evidence (2026-07-31):** traced direct production imports,
callers, results, errors, cancellation, and architecture guards for
`cmd/quill`, `internal/cli`, `internal/engine`, `internal/report`, and
`internal/architecture`. The current topology is the desired termination
state: no package move, new layer, public facade, or framework is justified.
The completed CLI-boundary refactor remains pinned under the exclusion above.
The review corrected the canonical dependency diagram, documented the actual
engine workflow surface, narrowed option contracts to their real seams,
made engine toolchain validity authoritative through report, and strengthened
high-level import enforcement.

**Why now:** this phase asks whether the repository is divided by durable
responsibility rather than by frameworks, technical layers, or historical edits.
It is the primary architectural review.

**Review paths:**

- `cmd/quill`, `internal/cli`, `internal/engine`, `internal/report`, and
  `internal/architecture`.
- `docs/architecture.md`, ADRs, and direct import graph evidence.

**Review questions:**

- Does each top-level package own a real concept with a single reason to change?
- Are `cmd/quill`, CLI, engine, and report responsibilities crisp and non-
  overlapping?
- Is `Engine` a deep private application facade, or has it become a grab bag of
  unrelated orchestration?
- Are CLI commands thin adapters that parse, choose streams, invoke engine, and
  apply public policy without second workflow paths?
- Does report convert engine results into explicit public views without pulling
  presentation into engine or Checks?
- Are architecture tests enforcing meaningful import/ownership seams?
- Does the filesystem layout make reader flow obvious, or do files and packages
  reveal historical migration shapes instead of current concepts?

**Exit criteria:** a package-by-package ownership map, a list of real seams,
and a clear conclusion on whether the current topology is the desired
termination state.

**Verdict:** `cmd/quill` is the shallow process shell; CLI is the sole inbound
adapter; engine is the cohesive private workflow facade; report is the
presentation adapter; and architecture is executable source policy. The
remaining package topology is capability-oriented and coherent.

**Re-review trigger:** a new entrypoint, public Go interface, service or
delivery mechanism; a CLI import of Profile, Pack, execution, installer,
toolchain, process, or filesystem orchestration; engine importing CLI/report;
report choosing execution, stream, or exit policy; or a material change to ADR
0001, 0004, 0005, or 0006.

### 2. Core conceptual model: Profile, Packs, Rules, Tools, and execution

**Status:** Complete.

**Verdict (2026-07-31):** The Profile -> Pack -> RuleDefinition/Template ->
Plan/Rule/Job -> Driver -> result path is single-owned and coherent. Profile owns
strict persisted-shape validation and generic containment; Pack owns catalogue
identity, selection, and Pack Policy semantics; `style` owns the closed
execution vocabulary; engine joins repository documents into a fresh operation
snapshot; execution dispatches the known Job families; report renders derived
truth. Shipped Pack runtime bindings and self-describing External Pack Jobs are
the two justified runtime seams. No new package, registry, framework, or
parallel External-Pack pipeline is warranted.

**Corrections and evidence:** Profile now accepts only table-shaped raw Pack
Policy and defers empty-object meaning to the owning Pack; Pack identities
reject blank and reserved values; repository-relative Profile paths validate the
deepest existing symlink-resolved ancestor; Rule-list text/JSON derives Pack
provenance; and Tools no longer advertise unsupported External provenance.
Strict Profile, Pack, execution, Driver, External Pack, metadata, report, and
CLI tests passed, as did `go test ./...`, `go test -race ./...`,
`golangci-lint run ./...`, `go build ./cmd/quill`, and list JSON/text smoke
tests. Focused review rechecked the reserved-Pack-ID correction.

**Re-review trigger:** a new Template or Job family; changed Check/Fix
cardinality; External Pack fixes, tools, file-set defaults, remote sources, or
protocol version; altered Pack-wide Tool or Policy lifecycle; changed identity,
containment, provenance, selection, or Runtime Binding semantics; a new
persisted Profile Policy shape; or another catalogue, compilation, or execution
path.

**Why now:** this is Quill's centre of gravity. It determines whether consumer
policy compiles coherently into reusable capabilities and executable work.

**Review paths:**

- `internal/style/*` and `internal/profile/*`.
- `internal/pack/*`, `internal/pack/shipped/**`, and
  `internal/pack/external/*`.
- `internal/execution/*`, `internal/execution/drivers/*`, and
  `internal/toolchain/*`.
- Relevant tests, Profile fixtures, and root `quill.toml`.

**Review questions:**

- Is there one coherent path from persisted Profile to Plan to Rule to Job to result?
- Do `style` and `profile` form an appropriately low dependency floor, or do
  higher-level concerns leak into them?
- Does Profile own consumer configuration while each Pack owns typed Pack
  policy, declarations, and runtime bindings?
- Does Pack composition remain explicit and deterministic without mutable
  global registration or a speculative framework?
- Are Rule, Tool, capability, Target, scope, file set, requirement, Check, and
  Fix relationships represented with the right cardinality and ownership?
- Are Driver interfaces justified by genuine variation across execution shapes,
  and do they hide complexity rather than mirror implementations?
- Is external Pack support a clean adapter seam with a deliberately smaller
  interface than its protocol implementation?

**Exit criteria:** a simple model diagram and a judgement on whether the
Profile/Pack/Rule/Tool execution model is conceptually minimal, complete, and
extensible only where real variation exists.

### 3. Major-module interface depth and local code quality

**Status:** Complete.

**Verdict (2026-07-31):** Quill's major modules expose narrow interfaces at
real ownership seams. The repository model, host adapters, execution adapters,
presentation packages, concrete Checks, and support code each retain local
complexity rather than recreating a second application layer. No generic
registry, one-use interface, forwarding alias, speculative constructor, or
line-count-driven package split is warranted.

**Corrections and evidence:** This review found three high-leverage local
defects, all corrected:

- `report` now accepts one `ToolchainResult` DTO, preserves engine-owned
  `AllValid`, and has text/JSON writers consume that same presentation input.
- `filewalk` now exposes one roots-first `CollectFiles` operation for regular
  files, with `CollectAllFiles` retaining its explicit non-binary contract.
  Go check selection and execution file-set collection use the same traversal
  seam.
- `workspace.ResolveRepoRelative` and Profile's mirrored resolver now enforce
  physical containment for missing paths below absolute and relative dangling
  symlinks, including a dangling link reached through a safe symlink alias.

Focused contract tests cover the new traversal and containment semantics.
The complete suite passed: `go test ./...`, `go test -race ./...`,
`golangci-lint run ./...`, and `go build ./cmd/quill`. An independent review
found and rechecked the dangling-symlink edge cases; no remaining defect was
found in the reviewed cutovers.

**Re-review trigger:** a new exported module seam; a new Driver, Pack runtime,
or report format; a second host implementation; a changed repository
containment or file-selection contract; or a helper that starts coordinating
more than its owning package.

**Why now:** after the overall architecture is accepted or revised, review each
major module for elegance: small interfaces, local complexity, natural names,
reader flow, and testable responsibilities.

**Review groups:**

- **Repository model:** `workspace`, `filewalk`, `styleguide`, `coverage`,
  `lockfile`, `markers`, and Profile path/compilation helpers.
- **Host adapters:** `process`, `installer`, `ecosystem`, and `toolchain`.
- **Execution model:** `execution`, generic Drivers, shipped Pack bindings,
  external Pack adapter, and engine operation methods.
- **Presentation model:** report views/renderers, CLI command adapters, and
  command entrypoint lifecycle.
- **Policy observations:** concrete `checks` packages and Pack-local codecs.
- **Support:** `testutil`, testdata, architecture tests, workflow, and build
  configuration.

**Review questions for every module:**

- Can a caller understand the interface without knowing implementation details?
- Are exported names precise domain nouns rather than mechanical or historical
  labels?
- Does a helper name a real concept, or should it be inlined?
- Is a file split an ownership boundary or merely a line-count reaction?
- Are constructors and options necessary at the module's true seam?
- Are there one-use interfaces, forwarding methods, aliases, registries, or
  abstractions without a second real implementation or variation point?
- Are errors and results structured at the owner that can make the correct
  decision?
- Do tests cross the same interface callers use, or do they depend on internal
  plumbing without defending an observable contract?

**Exit criteria:** a prioritised module-quality list that separates defects and
high-leverage simplifications from taste-only churn.

### 4. Core workflows and operation coherence

**Status:** Complete.

**Verdict (2026-07-31):** Repository operations retain one fresh
Profile/STYLE.md/Pack/Plan preparation path, while metadata-only Coverage,
Metadata, and Explain intentionally avoid executable work. Engine owns
operation-specific side effects and result values; report renders those values;
CLI alone selects streams and exit status. Init correctly remains separate
because it creates, rather than loads, a repository.

**Corrections and evidence:** Tool inspection now preserves parent-context
cancellation as an operation error instead of converting it into an invalid Tool
status. Check, Fix, Inspect, Install, and `IsInstalled` propagate that outcome.
Installation and archive resolution stop before beginning another Tool or
platform after cancellation; platform resolution is deterministic; and
context-aware atomic lock writing checks cancellation immediately before rename,
leaving the existing lockfile unchanged. Focused tests cover cancellation during
the first and final probes, empty inspection, post-first-tool cancellation,
sorted platform processing, and cancellation at the lock commit point. An
independent review found and rechecked both cancellation windows.

`prepareRun` still enriches Tool-only operations with Driver bindings, and
Metadata composes a complete catalogue twice. Both are low-cost, internal
simplifications without a current behavioural defect; this audit deliberately
does not add a second preparation path or a cache.

`go test ./...`, `go test -race ./...`, `golangci-lint run ./...`, and
`go build ./cmd/quill` passed. `quill coverage --format json
--repository-root .` completed successfully.

**Re-review trigger:** a new operation; changed cancellation, partial-side
effect, or lock-publication guarantee; a new Driver/Tool runtime family; altered
error/result ownership; or a second snapshot, catalogue, or preparation path.

**Why now:** this phase tests whether the conceptual model remains coherent when
Quill performs real operations rather than only when packages are read alone.

**Review paths:**

- `internal/engine/*`, `internal/execution/*`, `internal/report/*`,
  `internal/coverage/*`, and relevant CLI command files.
- Operation tests for check, fix, install, lock, coverage, metadata, list,
  explain, init, and external Packs.

**Review questions:**

- Does every engine operation share one preparation path without hidden
  duplicate catalogues, Profiles, or policy resolution?
- Are metadata-only workflows appropriately lighter than execution workflows?
- Do Check, Fix, Inspect, Install, Lock, Coverage, Metadata, Explain, and Init
  have crisp responsibilities and intentional overlap?
- Are result types explicit enough for report without leaking renderer or CLI
  concerns into engine?
- Is error classification made once at the correct level rather than translated
  repeatedly through every layer?
- Does cancellation travel through the whole workflow without creating special
  side paths or ambiguous partial-state semantics?

**Exit criteria:** an operation matrix naming each workflow's inputs,
preparation, side effects, result, error owner, rendering path, and tests.

### 5. Concrete Checks, fixes, and policy expression

**Status:** Complete.

**Verdict (2026-07-31):** The six shipped Check families remain Pack-qualified
observations. Every shipped Check, fixer, Profile Check, target command/check,
repository scan, and external Check dispatches explicitly to a runtime binding.
Pack-local codecs resolve and validate typed Pack Policy before Profile
compilation. Diagnostics retain their owning `style.Rule` identity through
engine and report, and sorted collection establishes deterministic traversal.

No cross-Pack policy abstraction, shared generic fixer, or Check registry is
warranted. The implementations differ in parsing, file selection, and
diagnostic semantics; common infrastructure would make those differences less
local rather than reduce real duplication. The focused independent review found
no correctness defect or evidence-backed change to make.

**Re-review trigger:** a new Pack/Check/fixer; a policy field that crosses
Pack ownership; changed diagnostic schema or ordering; new generated-file or
exception semantics; or a fix with partial-write behavior.

**Why after the core model:** individual rules should be judged against an
accepted model of Pack ownership, execution semantics, and diagnostics. They
are important but should not decide the architecture by accident.

**Review paths:**

- `internal/checks/{bash,golang,project,security,text,vocabulary}/**`.
- Pack-local binding and policy packages, Go scenario tests, and root
  rule-to-requirement bindings in `quill.toml`.

**Review questions:**

- Does each Check belong to its Pack and depend only on validated policy and
  execution values it actually needs?
- Are algorithms, parsing choices, and internal helper boundaries elegant and
  proportionate to the rule's value?
- Do diagnostics use a consistent model with stable rule identity, valid range,
  actionable wording, and deterministic order?
- Are false positives, false negatives, parser/type failures, generated files,
  exceptions, scope selection, and cross-file relationships handled at the
  correct layer?
- Are fixes safe, scope-confined, idempotent where appropriate, and explicit
  about partial failure?
- Are Pack policy codecs concise and language-neutral outside Pack-specific
  behaviour?

**Exit criteria:** a Pack-by-Pack quality assessment separating rule semantics,
shared algorithm opportunities, and unwanted cross-Pack coupling.

### 6. Operational trust boundaries and platform mechanics

**Status:** Complete - remediation required.

**Verdict (2026-07-31):** The review confirmed strict Profile/Pack containment,
direct argv execution, archive validation before extraction, checksum-before-
extraction, and Unix process-group cancellation. It also found six
evidence-backed defects:

- **Critical QUILL-TRUST-006:** repository cache PATH entries can execute a
  checked-out `go` or `npm` during bootstrap installation.
- **High QUILL-TRUST-001:** Go/npm install paths follow cache-directory
  symlinks outside the repository.
- **High QUILL-TRUST-003:** extension-based collection admits leaf symlinks and
  special files, reaching reads and mutating file-command tools.
- **Medium QUILL-TRUST-002:** archive destination pathname checks retain a
  symlink TOCTOU window before temporary-file creation and rename.
- **Medium QUILL-TRUST-005:** lockfile loading follows links and has no size
  bound.
- **Medium QUILL-TRUST-004:** Windows cancellation kills only the direct child,
  leaving descendants alive.

The first correction has an architecture recommendation: retain cache-first
PATH for normal Tool execution, but give bootstrap Go/npm a separate
repository/state-excluding PATH and resolve their absolute host executable only
after the installed-Tool early return. The remaining findings require their
own containment and platform remediation with hostile regression tests.

**Re-review trigger:** any remediation, new installer/runtime family, new
filesystem write or scanner, external Pack execution change, or Windows
lifecycle support.

**Why here:** these are essential correctness and security reviews, but their
implementation detail should be judged after the design explains what each
boundary is responsible for protecting.

**Review paths:**

- `internal/workspace/*`, `internal/profile/{path,path_resolution}.go`,
  `internal/filewalk/*`, `internal/engine/init.go`, and
  `internal/installer/destination.go`.
- `internal/installer/*`, `internal/lockfile/*`, `internal/process/*`,
  `internal/ecosystem/*`, `internal/toolchain/*`, and
  `internal/pack/external/*`.
- `SECURITY.md`, rejection-path tests, Unix/Windows lifecycle tests, and
  archive/download fixtures.

**Check filesystem boundaries:**

- canonical-root and repository-relative resolution;
- NUL, absolute, drive, parent, symlink, and missing-descendant containment;
- regular-file selection, exclusion, generated-file treatment, and traversal
  bounds;
- no-clobber, rollback, parent handling, permissions, staging, and cleanup for
  init, lockfile, cache, and tool writes.

**Check network, archive, and process boundaries:**

- HTTPS and redirect policy, download limits, checksum verification order,
  archive path/type/link/size validation, and atomic replacement;
- direct argv execution, executable identity, environment inheritance versus
  isolation, output limits, timeouts, cancellation, process-tree cleanup, and
  platform differences;
- external Pack manifest strictness, executable containment, protocol limits,
  JSON Lines validation, diagnostic-range validation, and stderr handling.

**Exit criteria:** a boundary table with attacker-controlled values, validation
points, trusted models, remaining TOCTOU assumptions, and hostile regression
coverage.

### 7. Public integration, protocol, and release quality

**Status:** Complete - remediation required.

**Verdict (2026-07-31):** CI/release permissions, native archive/checksum
verification, release smoke tests, and the CLI separation are coherent.
Four public mismatches require correction: the protocol classifies completed
per-Rule execution errors as command-level `operation_failed`; documented
version strings omit the executable's `v` prefix and `(devel)` fallback; the
generated Init instruction uses invalid `quill list`; and Init says both policy
files exist when either file prevents creation. These are user-visible,
repository-local corrections.

**Re-review trigger:** schema/version/exit/stream change, new command, altered
Init template, or CI/release workflow modification.

**Why late:** public delivery must reflect a workflow model already proven
coherent. This is protocol and product review, not a redesign of the pinned
CLI boundary.

**Review paths:**

- `cmd/quill/*`, `internal/cli/*`, `internal/report/*`, and
  `internal/architecture/*`.
- `docs/cli-protocol.md`, CLI help fixtures, grammar tests,
  `machine_protocol_test.go`, command-entrypoint tests, and version tests.
- `.github/workflows/{ci,release}.yml`, `Makefile`, `go.mod`, and `go.sum`.

**Review questions:**

- Does the CLI expose the product model without leaking Kong or engine details?
- Are help, grammar, streams, JSON envelopes, text output, errors, exits,
  cancellation, and version metadata deterministic and compatible?
- Do report DTOs define a stable public view rather than serialising private
  engine types?
- Do CI and release workflows prove the actual platform/archive/version claims?
- Are checksums, archive layout, native smoke tests, action permissions, and
  rerun behaviour proportionate to the distribution contract?

**Exit criteria:** protocol and release evidence that is behavioural, platform-
aware, and independent of internal implementation shape.

### 8. Tests, self-hosting, documentation, and final polish

**Status:** Complete - remediation required.

**Verdict (2026-07-31):** Architecture guards, Profile round trips, scenario
tests, help goldens, and current lint/test/build gates give a credible
baseline. Assurance is incomplete for immutable public schemas and risk
boundaries: contract tests accept missing documented JSON fields; the External
Pack wire test does not assert `scope`; Windows tests and race detection are
not CI gates; tag-only release shell logic lacks PR-time action/shell checking;
and focused regression coverage is absent for the trust defects above.

**Re-review trigger:** any public schema, External Pack request, platform
support, CI/release gate, or trust-boundary correction.

**Why last:** these prove, communicate, and maintain the earlier conclusions.
They cannot substitute for direct review of the model and implementation.

**Review paths:**

- `internal/testutil/**`, all `_test.go` files, testdata, and scenario tests.
- `STYLE.md`, `quill.toml`, `quill.lock`, `.golangci.yml`, `Makefile`, root
  documents, and `docs/`.

**Review questions:**

- Do tests defend observable contracts, error paths, architecture constraints,
  trust-boundary rejection, and platform behaviour rather than source shape?
- Are fixtures deterministic, isolated, and safe with race detection?
- Does self-hosted Quill policy catch meaningful quality regressions without
  becoming circular reassurance or forcing arbitrary code style?
- Does every durable document have one owner, use the same vocabulary, and
  describe the actual code and product contracts?
- Are names, comments, errors, helpers, file order, and formatting coherent
  enough that future contributors can navigate the system naturally?

**Completion criteria:**

- Every production package and supported end-to-end flow is accounted for by a
  primary phase.
- Every design finding has a proposed coherent termination state, a justified
  scope, and a disposition in its canonical owner.
- Every correctness/security finding has focused regression evidence.
- Relevant narrow tests and repository-wide quality/release gates have observed
  evidence when their phases complete.
- No phase is marked complete merely because self-lint passes.

## Cross-cutting review passes

Apply these passes throughout the programme:

- **Naming and terminology:** use precise domain nouns; eliminate historical,
  mechanical, or duplicate names only when the replacement improves locality.
- **Error and cancellation semantics:** preserve causes, classify them once,
  release resources once, and make partial-state behaviour explicit.
- **Determinism:** define ordering rather than relying on map, filesystem, or
  tool-output order.
- **Compatibility:** version and migrate JSON, repository formats, error codes,
  stream behaviour, exit semantics, help, and release artefacts deliberately.
- **Performance and bounds:** keep file walking, parsing, output capture,
  archive extraction, and protocol parsing bounded by their actual risk.
- **Platform behaviour:** make Unix and Windows differences explicit and test
  them; build success is not proof of lifecycle equivalence.
- **Documentation drift:** correct the canonical owner when code and docs
  disagree; do not duplicate stale summaries here.

## Finding disposition and programme maintenance

Route review outcomes deliberately:

- **Small bounded repository-local implementation task:** `TODO.md` or the
  tracked work item that replaces it.
- **MVP scope, acceptance criteria, or delivery dependency:** `docs/mvp.md`
  and/or `docs/roadmap.md`.
- **Durable package-boundary or ownership decision:** a new or superseding ADR.
- **Current package ownership or runtime-flow correction:**
  `docs/architecture.md`.
- **JSON, stream, exit, cancellation, or CLI compatibility change:**
  `docs/cli-protocol.md` and product/release documentation as needed.
- **Security invariant, threat-model, or reporting correction:** code and tests
  plus `SECURITY.md` where the model changes.
- **Concrete implementation defect:** the owning package, focused regression
  test, and relevant canonical documentation.

Update this roadmap only when review scope/order changes, a phase begins,
completes, blocks, or reopens, a material package/flow moves, or the latest
baseline/evidence changes. Replace stale phase records rather than accumulating
historical transcripts. Reopen a completed phase when its inputs, outputs,
trust boundary, public contract, or accepted ADR materially changes.
