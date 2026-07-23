# Project Style Guide

## 0. Governance and enforcement

### 0.1 Normative language

Required rules:

<!-- style:
id=0.1.rfc-keywords
mode=review_only
reason=Editorial policy; enforce in review, not lint.
-->
* The keywords "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are normative.

Priority order for rule conflicts:

<!-- style:
id=0.1.security-first
mode=review_only
reason=Priority ordering is architectural judgement.
-->
1. Security and correctness rules
<!-- style:
id=0.1.architecture-second
mode=review_only
reason=Priority ordering is architectural judgement.
-->
1. Architecture rules
<!-- style:
id=0.1.language-third
mode=review_only
reason=Priority ordering is editorial judgement.
-->
1. Language-specific style rules
<!-- style:
id=0.1.formatting-fourth
mode=review_only
reason=Priority ordering is editorial judgement.
-->
1. Formatting rules
<!-- style:
id=0.1.recommendations-last
mode=review_only
reason=Priority ordering is editorial judgement.
-->
1. Recommendation rules

<!-- style:
id=0.1.local-exception
mode=review_only
reason=Exception scope and compatibility justification remain contextual.
-->
* When a rule conflicts with generated code, vendored code, or a documented compatibility
  requirement, the exception MUST be local, explicit, and justified inline.

### 0.2 Excluded files and directories

Unless a section states otherwise, the following are excluded from style enforcement:

<!-- style: id=0.2.exclude-vendor -->
* `vendor/`
<!-- style: id=0.2.exclude-third-party -->
* `third_party/`
<!-- style: id=0.2.exclude-testdata -->
* `testdata/`
<!-- style: id=0.2.exclude-generated -->
* Generated files that include `DO NOT EDIT.`
<!-- style:
id=0.2.exclude-external-fixtures
mode=review_only
reason=Whether a file is an external fixture depends on repository context.
-->
* Fixture files copied from external sources
<!-- style:
id=0.2.exclude-crypto-vectors
mode=review_only
reason=Whether a file is a cryptographic test vector depends on repository context.
-->
* Cryptographic test vectors
<!-- style:
id=0.2.exclude-migrations
mode=review_only
reason=Migration-snapshot classification depends on repository context.
-->
* Migration snapshots

<!-- style: id=0.2.no-committed-secrets -->
* These exclusions do not permit unsafe secret handling. Secrets MUST NOT be committed anywhere.

### 0.3 Exception markers

Exception markers MUST use this exact form:

<!-- style: id=0.3.exception-basic -->
* `style: allow-<rule>`
<!-- style: id=0.3.exception-with-reason -->
* `style: allow-<rule> because: <short reason>`

Rules for exception markers:

<!-- style:
id=0.3.exception-local
mode=review_only
reason=Whether a same-line marker is practical depends on local readability.
-->
* Keep markers on the same line as the exception when practical.
<!-- style: id=0.3.exception-ascii -->
* Keep markers ASCII only.
<!-- style:
id=0.3.exception-specific
mode=review_only
reason=Reason specificity still needs human judgement.
-->
* Keep reasons short and specific.
<!-- style:
id=0.3.exception-narrow
mode=review_only
reason=Narrowest-possible exceptions still need contextual review.
-->
* Use the narrowest possible exception.
<!-- style:
id=0.3.exception-filewide-justified
mode=review_only
reason=Compatibility justification for file-wide exceptions is review-driven.
-->
* File-wide exceptions require a documented compatibility reason.

Examples:

* `# style: allow-long-line because: shell output must remain copy-pasteable`
* `// style: allow-non-ascii because: protocol sample requires UTF-8`

### 0.4 Toolchain and formatter versions

The repository MUST pin versions for:

<!-- style: id=0.4.pin-go -->
* Go
<!-- style: id=0.4.pin-golangci-lint -->
* `golangci-lint`
<!-- style: id=0.4.pin-goimports -->
* `goimports`
<!-- style: id=0.4.pin-shfmt -->
* `shfmt`
<!-- style: id=0.4.pin-shellcheck -->
* `shellcheck`
<!-- style: id=0.4.pin-markdownlint -->
* `markdownlint`
<!-- style: id=0.4.pin-spelling-checker -->
* The spelling checker

<!-- style: id=0.4.use-pinned-toolchain -->
* CI and local style scripts MUST use the pinned versions.

### 0.5 Enforcement levels

<!-- style: id=0.5.required-fails-ci -->
* Required rules MUST fail CI.
<!-- style: id=0.5.recommendations-report-only -->
* Recommendation rules MUST be reported in CI but MUST NOT fail CI by default.
<!-- style: id=0.5.experimental-opt-in -->
* Experimental rules MAY run only in an opt-in strict mode.

The repository MUST provide:

<!-- style: id=0.5.make-lint -->
* `make lint`
<!-- style: id=0.5.make-lint-required -->
* `make lint-required`
<!-- style: id=0.5.make-lint-fix -->
* `make lint-fix` where safe auto-fixes exist
<!-- style: id=0.5.make-style-install -->
* `make style-install`
<!-- style: id=0.5.make-style-doctor -->
* `make style-doctor`
<!-- style: id=0.5.make-style-coverage -->
* `make style-coverage`

### 0.6 Style profile and repository policy

<!-- style:
id=0.6.style-guide-human-source
mode=review_only
reason=Source-of-truth ownership is a governance rule rather than a lintable syntax rule.
-->
* `STYLE.md` is the human-readable normative source of repository style policy.
<!-- style:
id=0.6.style-profile-machine-source
mode=review_only
reason=The machine-profile role is architectural documentation rather than a syntax rule.
-->
* `quill.toml` is the machine-readable repository profile consumed by Quill.
<!-- style:
id=0.6.coverage-from-profile-bindings
mode=review_only
reason=Coverage derivation spans profile data and guide metadata,
so semantic alignment still needs review.
-->
* Automated coverage is derived from active `quill.toml` rule bindings plus `STYLE.md` requirement
  metadata.
<!-- style:
id=0.6.keep-guide-and-profile-aligned
mode=review_only
reason=Semantic alignment between prose policy and machine profile still needs human judgement.
-->
* Repository-specific policy changes MUST keep `STYLE.md` and `quill.toml` aligned.
<!-- style:
id=0.6.keep-style-platform-language-agnostic
mode=review_only
reason=Platform modularity is an architectural design property, not a syntax rule.
-->
* Generic style-platform execution, file selection, and reporting SHOULD avoid hardcoding individual
  programming or document languages.
<!-- style:
id=0.6.isolate-language-specific-rules
mode=review_only
reason=Deciding whether a concern belongs in a generic runner or language module
still needs design judgement.
-->
* Language-specific analysis SHOULD live in explicit language modules or profile data, not in
  generic runner dispatch.
<!-- style:
id=0.6.rule-packs-own-checker-capabilities
mode=review_only
reason=Rule-pack boundaries are architectural structure, not a syntax-only rule.
-->
* Checker capabilities SHOULD be grouped into explicit rule packs rather than a single central
  registry.
<!-- style:
id=0.6.requirement-ids-are-project-policy
mode=review_only
reason=Requirement ownership is partly enforced by tests and partly by review.
-->
* Requirement IDs belong to `STYLE.md` and `quill.toml`; implementation code SHOULD use checker-
  owned diagnostic codes instead.

## 1. Repository-wide conventions

### 1.1 Readability and consistency principles

Scope:

* Applies to repository-wide judgement calls across code, comments, documentation, and scripts.

Required rules:

<!-- style:
id=1.1.correctness-first
mode=review_only
reason=This is a review principle, not a mechanical rule.
-->
* Prefer correctness over convenience.
<!-- style:
id=1.1.clarity-over-cleverness
mode=review_only
reason=Clarity at the design level still needs human review.
-->
* Prefer clarity over cleverness.
<!-- style:
id=1.1.simplicity-over-cleverness
mode=review_only
reason=Simplicity at the design level still needs human review.
-->
* Prefer simplicity over clever abstraction.
<!-- style:
id=1.1.prefer-standard-library-helpers
mode=review_only
reason=Whether a helper adds useful domain meaning still needs review.
-->
* Prefer standard library helpers such as `slices.Contains`, `slices.ContainsFunc`, `maps.Keys`, and
  `slices.Sort` when they express the operation directly. Keep a local helper only when it adds
  domain meaning or centralises non-trivial policy.
<!-- style:
id=1.1.consistency-with-surrounding-code
mode=review_only
reason=Consistency with nearby code is contextual rather than purely syntactic.
-->
* Prefer consistency with surrounding code unless a stronger repository rule applies.

Recommendation rules:

<!-- style:
id=1.1.follow-local-convention
mode=review_only
reason=Local convention is contextual rather than purely syntactic.
-->
* When multiple acceptable forms exist, prefer the form already established in the surrounding file
  or package unless a documented repository rule says otherwise.

### 1.2 Line length

Scope:

* Applies to code, comments, documentation, and user-facing strings.

Required rules:

<!-- style: id=1.2.max-line-length -->
* Maximum 100 characters per line.
<!-- style:
id=1.2.wrap-thoughtfully
mode=review_only
reason=Wrapping quality is an editorial judgement.
-->
* Prefer wrapping thoughtfully over horizontal scrolling.
<!-- style:
id=1.2.markdown-reference-links
mode=review_only
reason=Choosing reference links and line placement is an editorial judgement.
-->
* In Markdown, prefer reference links for long inline destinations, putting the destination on a
  following source line when doing so keeps every line within 100 columns.
<!-- style: id=1.2.tabs-count-four -->
* Tabs count as 4 columns for line-length checks.

Exceptions:

<!-- style: id=1.2.generated-exception -->
* Generated code, for example files with `// Code generated ... DO NOT EDIT.`
<!-- style: id=1.2.vendor-exception -->
* Vendored code
<!-- style: id=1.2.machine-maintained-exception -->
* Machine-maintained checksum, lock, and local tool-state files, for example `go.sum`, `package-
  lock.json`, and `.cache/`
<!-- style:
id=1.2.raw-string-exception
mode=review_only
reason=Raw-string readability is contextual.
-->
* Long raw string literals where wrapping harms readability, for example SQL, PEM blocks, and test
  vectors
<!-- style:
id=1.2.localise-long-content
mode=review_only
reason=Content locality is a readability judgement.
-->
* In those cases, keep the content in `const` blocks and avoid long lines elsewhere.
<!-- style:
id=1.2.shell-output-exception
mode=review_only
reason=Shell-output clarity is contextual.
-->
* Long user-facing shell output lines are allowed only when wrapping would materially reduce
  clarity.
<!-- style: id=1.2.shell-long-line-marker -->
* Any shell-line exception MUST include an inline `# style: allow-long-line` marker.
<!-- style: id=1.2.markdown-reference-destination-exception -->
* A Markdown source line may exceed 100 columns only when that line holds a valid link reference
  definition whose destination is one exact lowercase `http://` or `https://` URL whose own text
  is strictly longer than 100 columns, and whose non-destination remainder (label, delimiters,
  spacing, title, and other content) is at most 100 columns. Inline links, bare URLs, autolinks,
  code and HTML URLs, invalid definitions, non-HTTP(S) destinations, and excess label or title
  text MUST still fail.

Automation guidance:

* Enforce line length through the profile-owned `line_length` file set.
* Bash line-length checks SHOULD treat tabs as 4 columns.
* Shell-line exceptions SHOULD require an explicit inline marker.
* Quill is the sole Markdown line-length authority; markdownlint `MD013` is disabled to avoid a
  competing check, so only the reference-definition destination exception above applies.

### 1.3 Language

Scope:

* Applies to documentation, comments, user-facing CLI output, log messages, and error context
  strings.

Required rules:

<!-- style: id=1.3.australian-english -->
* Use Australian English everywhere in the codebase.
<!-- style: id=1.3.configured-spelling-wins -->
* When multiple accepted English spellings exist, prefer the form required by the configured
  spelling checker and the existing codebase convention.

Automation guidance:

* Use tooling configured for UK English as a practical proxy for Australian English.
* Prefer automated spelling checks over manually maintained dictionaries.

Examples:

* `initialise`, `colour`, `behaviour`, `optimise`, `centre`, `organise`, `authorisation`

### 1.4 Character set

Scope:

* Applies to all repository files unless non-ASCII is strictly necessary.

Required rules:

<!-- style: id=1.4.ascii-default -->
* Use ASCII across the repository by default.
<!-- style: id=1.4.source-files-ascii -->
* Source files MUST be ASCII by default.
<!-- style: id=1.4.docs-and-strings-ascii -->
* Comments, documentation, and user-facing strings MUST be ASCII by default.
<!-- style:
id=1.4.non-ascii-minimal
mode=review_only
reason=Minimal and intentional Unicode usage is still a readability judgement.
-->
* If non-ASCII is required, keep its use minimal and intentional.
<!-- style: id=1.4.non-ascii-marker -->
* Any non-ASCII exception MUST include an inline `style: allow-non-ascii` marker.

Automation guidance:

* Use automated ASCII checks across both app code and tooling code.
* Enforce non-ASCII exceptions via an explicit inline marker.

Examples:

* Prefer `centre` over typographic variants that require Unicode punctuation.
* Prefer plain `'` and `"` over curly quotes.
* Keep shell output ASCII unless Unicode meaningfully improves clarity.

### 1.5 Section headers

Scope:

* Applies to Go and Bash source files.

Required rules:

<!-- style: id=1.5.use-section-headers -->
* Use block-comment section headers to reduce navigation cost between distinct conceptual groups
  within a file. Do not add headers as ritual.
<!-- style: id=1.5.headers-over-100-lines -->
* Use headers when the file has 100 or more lines.
<!-- style:
id=1.5.short-files-can-skip-headers
mode=review_only
reason=Need-for-headers still depends on file complexity.
-->
* Files under 80 lines or with a single responsibility typically do not need headers. Files from 80
  to 99 lines remain a judgement call: use headers when distinct local concepts would otherwise blur
  together, and omit them when the file reads as one uninterrupted narrative.
<!-- style: id=1.5.header-width -->
* Section header lines MUST be exactly 100 characters long.
<!-- style: id=1.5.header-centred -->
* Header text MUST be centred.
<!-- style: id=1.5.left-fill-precedence -->
* If the fill count cannot be perfectly equal, the left side MUST have exactly one more fill
  character than the right.
<!-- style: id=1.5.header-tabs-count-four -->
* Tabs count as 4 columns for width checks.
<!-- style:
id=1.5.shorten-long-header-text
mode=review_only
reason=Header wording is an editorial judgement.
-->
* If a section name is too long to fit, shorten the section name.

Recommendation rules:

<!-- style:
id=1.5.short-files-avoid-headers
-->
* Files under 80 lines SHOULD normally have no section headers.
<!-- style:
id=1.5.header-overuse-split-file
-->
* Treat 7 or more section headers in one file as a sign that the file should probably be split or
  simplified.
<!-- style:
id=1.5.structural-headers-large-mixed-files
mode=review_only
reason=Whether a structural label improves navigation depends on the file's role and mix of content.
-->
* Structural labels such as `Types`, `Constants`, and `Helpers` are acceptable only in larger mixed
  files where they improve navigation, not in tiny single-purpose files.
<!-- style:
id=1.5.prefer-specific-section-names
-->
* Prefer specific section names over generic labels such as `Check`, `Checks`, `Misc`, or `Other`,
  unless the label is already a repository-standard structural section and the file is large enough
  to justify it.
<!-- style:
id=1.5.match-heading-grammar-by-file-role
mode=review_only
reason=Heading-grammar consistency still depends on how file roles are grouped across the repo.
-->
* Files that serve the same role SHOULD reuse the same section-heading grammar.

Automation guidance:

* Use automated section-header checks for length, centring, left-side precedence, missing headers in
  large files, generic names, and header density. Header-density findings are recommendations.

Formats:

```go
/* ------------------------------------------- Config ------------------------------------------- */
```

```bash
# ------------------------------------------- Config --------------------------------------------
```

### 1.6 Magic values

Scope:

* Applies to Go and Bash code.

Required rules:

<!-- style: id=1.6.no-unexplained-numbers -->
* Do not embed unexplained numeric literals directly in logic.
<!-- style:
id=1.6.named-domain-constants
mode=review_only
reason=Whether a value is domain-significant is contextual.
-->
* Use named constants for protocol sizes, limits, timeouts, permissions, retry counts, and similar
  values.
<!-- style:
id=1.6.descriptive-constant-names
mode=review_only
reason=Name descriptiveness still needs human judgement.
-->
* Keep constant names descriptive and domain-specific.
<!-- style:
id=1.6.trivial-values-allowed
mode=review_only
reason=Obviousness of trivial literals is contextual.
-->
* Trivial values (`0`, `1`, `-1`) are allowed when their meaning is immediately obvious.

Recommendation rules:

<!-- style: id=1.6.bash-magic-values-recommendation -->
* In Bash scripts, magic-value findings are recommendation-level by default.

Automation guidance:

* Required enforcement for Go uses `golangci-lint` with `mnd`.
* Bash magic-value checks are reported as recommendations in `make lint`.

Correct:

```go
const receiveTimeoutSeconds = 10
const filePermOwnerOnly = 0o600

ctx, cancel := context.WithTimeout(context.Background(), receiveTimeoutSeconds*time.Second)
if err := os.MkdirAll(path, filePermOwnerOnly); err != nil {
    return err
}
```

```bash
MAX_RETRIES=5
for ((attempt = 0; attempt < MAX_RETRIES; attempt++)); do
	:
done
```

Incorrect:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
_ = os.MkdirAll(path, 0o600)
```

```bash
for ((attempt = 0; attempt < 5; attempt++)); do
	:
done
```

<!-- style: id=1.6.constants-use-mixedcaps -->
* Constants MUST use MixedCaps: `MaxRetries`, `DefaultTimeout`, `SectionSlug`. `ALL_CAPS` constants
  such as `MAX_RETRIES` are a C and Java convention and are forbidden in Go.

### 1.7 Call-site readability

Scope:

* Applies to function and method calls in Go and Bash code.

Required rules:

<!-- style:
id=1.7.no-mystery-literals-at-call-sites
mode=review_only
reason=Literal readability still needs API-level context.
-->
* Avoid passing unexplained boolean or numeric literals to functions.
<!-- style:
id=1.7.prefer-named-call-site-context
mode=review_only
reason=Call-site readability still depends on API and domain context.
-->
* Prefer named constants, option types, or brief inline parameter comments where needed.

### 1.8 General readability recommendations

Scope:

* Applies to readability patterns that are useful but not mandatory by default.

Recommendation rules:

<!-- style: id=1.8.blank-line-between-guards -->
* Prefer a blank line between consecutive guard clauses that each return early.
<!-- style:
id=1.8.separate-logical-steps
mode=review_only
reason=Logical-step boundaries are contextual.
-->
* Use blank lines to separate distinct logical steps in long functions.
<!-- style:
id=1.8.split-dense-blocks
mode=review_only
reason=Code density is a readability judgement.
-->
* Keep dense blocks readable by splitting validation, transformation, and side effects.
<!-- style: id=1.8.blank-lines-between-switch-cases -->
* In non-trivial `switch` statements, prefer a blank line between consecutive `case` blocks.
<!-- style:
id=1.8.prefer-guard-clauses
mode=review_only
reason=Control-flow preference is contextual.
-->
* Prefer guard clauses and early returns to deep nesting.
<!-- style:
id=1.8.limit-variable-scope
mode=review_only
reason=Ideal variable scope depends on readability and mutation flow.
-->
* Limit variable scope to the smallest practical block.

Automation guidance:

* Recommendation checks SHOULD be reported in `make lint`.
* Recommendation checks SHOULD NOT fail `lint-required` unless strict mode is selected.

### 1.9 Taxonomy, naming systems, and organisation

Scope:

* Applies to repository terminology, command surfaces, directories, files, and repeated helper
  families.

Required rules:

<!-- style:
id=1.9.one-term-per-concept
mode=review_only
reason=Concept boundaries and overlapping terms still need repository-level judgement.
-->
* Use one repository term for one concept.
<!-- style:
id=1.9.synonyms-require-difference
mode=review_only
reason=Whether two terms express a real semantic difference still needs design judgement.
-->
* If two terms are both kept, they MUST signal a real semantic difference.
<!-- style:
id=1.9.organise-by-one-axis
mode=review_only
reason=Primary directory and file grouping axes still need repository-level design judgement.
-->
* Organise sibling files and directories by one primary axis at a time, for example language,
  concern, or lifecycle stage.
<!-- style:
id=1.9.remove-dead-aliases
mode=review_only
reason=Whether a wrapper still carries distinct behaviour needs maintenance judgement.
-->
* Remove aliases, wrappers, and compatibility shims once they no longer provide distinct behaviour.
<!-- style:
id=1.9.similar-families-share-order
mode=review_only
reason=Comparable file families still need contextual review before enforcing a fixed order.
-->
* Similar file families MUST use a consistent internal order where practical.

Recommendation rules:

<!-- style:
id=1.9.prefer-family-consistency
mode=review_only
reason=Family-wide consistency versus local optimisation is a repository-level tradeoff.
-->
* Prefer family-wide consistency over isolated local renames.
<!-- style:
id=1.9.role-based-file-names
mode=review_only
reason=Role-versus-implementation naming remains contextual.
-->
* Prefer file and directory names that describe role or concern, not incidental implementation
  detail.
<!-- style:
id=1.9.avoid-mixed-taxonomy
mode=review_only
reason=Mixed grouping axes can still be justified in limited cases.
-->
* Avoid sibling sets that mix naming axes, for example language names beside behavioural categories,
  unless the distinction is explicit and justified.
<!-- style:
id=1.9.stable-naming-grammars
mode=review_only
reason=Naming grammar consistency still needs subsystem-aware review.
-->
* Prefer stable naming grammars for commands, reports, rule files, and helper families.

Examples:

* Prefer one style command family such as `style` unless a separate `lint` concept is intentionally
  preserved.
* Prefer role-specific section titles such as `Check Command`, `Check Output`, or `Naming Rules`
  over generic labels reused for unrelated file roles.

## 2. Architecture and security boundaries

### 2.1 Architecture

Scope:

* Keep shared domain and protocol logic independent from client, relay, CLI, filesystem, and network
  transport concerns.
* Treat the client and relay as separate hexagons that share only `internal/core` and
  `internal/relaywire`.
* Prefer explicit repository dependency rules over generic architecture slogans.

Directory responsibilities:

* `internal/core` Shared domain logic, protocol code, cryptography helpers, and relay boundary
  models.

  Rules:
  <!-- style: id=2.1.core-no-app-imports -->
  * MUST NOT import from `internal/client`, `internal/relay`, or `internal/relaywire`.
  <!-- style:
  id=2.1.core-stdlib-first
  mode=review_only
  reason=Dependency necessity is an architectural judgement.
  -->
  * Prefer the standard library only.
  <!-- style:
  id=2.1.core-third-party-essential
  mode=review_only
  reason=Dependency necessity is an architectural judgement.
  -->
  * Third-party dependencies are allowed only when essential for correctness or security.
  <!-- style:
  id=2.1.core-wrap-third-party
  mode=review_only
  reason=Wrapper boundaries are architectural.
  -->
  * Where practical, wrap third-party dependencies behind ports.
  <!-- style:
  id=2.1.core-allow-x-modules
  mode=review_only
  reason=Module choice is architectural.
  -->
  * `golang.org/x/...` modules are acceptable when appropriate.
  <!-- style:
  id=2.1.core-no-adapter-deps
  mode=review_only
  reason=Documented domain exceptions still need architectural review.
  -->
  * MUST NOT import transport, storage-driver, logging-framework, CLI, or network-client packages
    directly unless the dependency is required for a documented domain concern.

* `internal/relaywire` Shared relay HTTP payload models used by multiple hexagons but not part of
  the domain model.

  Rules:
  <!-- style: id=2.1.relaywire-no-app-imports -->
  * MUST NOT import from `internal/client` or `internal/relay`.
  <!-- style:
  id=2.1.relaywire-depend-on-core
  mode=review_only
  reason=Shared-package dependency choice is architectural.
  -->
  * MAY depend on `internal/core` and sibling shared packages.
  <!-- style:
  id=2.1.relaywire-shared-technical-focus
  mode=review_only
  reason=Package focus is architectural.
  -->
  * SHOULD stay focused on narrowly scoped technical concerns that are shared across hexagons.

* `internal/client/application/port` Client-side application ports.

* `internal/client/application/service` Client-side application services and orchestration.

* `internal/client/adapters/inbound` Client-side inbound adapters such as the CLI.

* `internal/client/adapters/outbound` Client-side outbound adapters such as filestore persistence
  and relay HTTP transport.

* `internal/client/bootstrap` Client composition root. This is the only client package that wires
  application code and adapters together.

* `internal/relay/application/port` Relay-side application ports.

* `internal/relay/application/service` Relay-side application services and policy orchestration.

* `internal/relay/adapters/inbound` Relay-side inbound adapters such as HTTP handlers.

* `internal/relay/adapters/outbound` Relay-side outbound adapters such as the in-memory queue store.

* `internal/relay/bootstrap` Relay composition root. This is the only relay package that wires
  application code and adapters together.

Dependency rules:

<!-- style: id=2.1.core-import-boundary -->
* Code in `internal/core` MUST NOT import from `internal/client`, `internal/relay`, or
  `internal/relaywire`.
<!-- style: id=2.1.port-dependencies -->
* Code in `internal/*/application/port` MUST depend only on `internal/core` and sibling
  `application/port` packages.
<!-- style: id=2.1.service-dependencies -->
* Code in `internal/*/application/service` MUST depend only on `internal/core`, sibling
  `application/port` packages, and sibling `application/service` packages.
<!-- style: id=2.1.inbound-dependencies -->
* Code in `internal/*/adapters/inbound` MUST stay inside its own application boundary and MAY depend
  on its application's `application/*` packages, `internal/core`, and `internal/relaywire`.
<!-- style: id=2.1.outbound-dependencies -->
* Code in `internal/*/adapters/outbound` MUST stay inside its own application boundary and MAY
  depend on its application's `application/port`, `internal/core`, and `internal/relaywire`.
<!-- style: id=2.1.bootstrap-dependencies -->
* Code in `internal/*/bootstrap` MAY depend on its own application's `application/*` packages, its
  own adapters, `internal/core`, and `internal/relaywire`.
<!-- style: id=2.1.cmd-ciphera-import -->
* `cmd/ciphera` MUST import only `internal/client/adapters/inbound/cli` from internal packages.
<!-- style: id=2.1.cmd-relay-import -->
* `cmd/relay` MUST import only `internal/relay/bootstrap` from internal packages.
<!-- style: id=2.1.shared-dtos-live-in-relaywire -->
* Shared relay DTOs that both applications need MUST live in `internal/relaywire`, not in a client-
  or relay-owned adapter package.
<!-- style:
id=2.1.env-reads-at-edge
mode=review_only
reason=Environment-read placement is architectural.
-->
* Environment reads SHOULD occur in inbound adapters or bootstrap packages, not in `internal/core`
  or `internal/*/application/service`.

Automation guidance:

* Use architecture import checks to enforce dependency boundaries.

### 2.2 Logging and observability

Scope:

* Applies to log messages and logging calls in Go and Bash code.

Required rules:

<!-- style: id=2.2.structured-logs -->
* Logs MUST use structured fields where the logger supports structured output.
<!-- style: id=2.2.stable-log-keys -->
* Log keys MUST be stable, lower-case, and ASCII.
<!-- style: id=2.2.no-secrets-in-logs -->
* Logs MUST NOT include secrets, tokens, private keys, passphrases, or raw credentials.
<!-- style:
id=2.2.actionable-error-logs
mode=review_only
reason=Actionability is a quality judgement rather than a syntax check.
-->
* Error logs SHOULD include actionable context but MUST avoid repeating the same full context at
  every abstraction layer.

Automation guidance:

* Use secret-scanning and logging checks where practical.

### 2.3 Secrets and sensitive data

Scope:

* Applies to source code, tests, fixtures, logs, comments, examples, and documentation.

Required rules:

<!-- style: id=2.3.no-secrets-in-repo -->
* Secrets MUST NOT appear in source, tests, fixtures, logs, comments, examples, or error strings.
<!-- style:
id=2.3.constant-time-comparison
mode=review_only
reason=Whether comparison is security-sensitive depends on domain context.
-->
* Security-sensitive equality checks MUST use constant-time comparison where appropriate.
<!-- style:
id=2.3.least-privilege-permissions
mode=review_only
reason=Least-privilege defaults still depend on the surrounding file contract.
-->
* File permissions MUST be least-privilege by default.
<!-- style:
id=2.3.normalise-untrusted-paths
mode=review_only
reason=Path-normalisation requirements depend on cross-function taint flow.
-->
* Paths derived from untrusted input MUST be validated and normalised before file operations.

### 2.4 Cryptography and randomness

Scope:

* Applies to security-sensitive tokens, keys, nonces, protocol state, and cryptographic code.

Required rules:

<!-- style: id=2.4.crypto-rand -->
* Use `crypto/rand` for keys, tokens, nonces, and other secret material.
<!-- style: id=2.4.no-deprecated-crypto -->
* Deprecated cryptographic algorithms MUST NOT be introduced without a documented compatibility
  requirement.
<!-- style:
id=2.4.stdlib-crypto-default
mode=review_only
reason=Default-library preference still needs human review when exceptions exist.
-->
* Cryptographic primitives SHOULD come from the Go standard library unless an approved exception
  exists.
<!-- style:
id=2.4.document-third-party-crypto
mode=review_only
reason=Documentation expectations are better enforced in review.
-->
* Third-party cryptographic dependencies MUST be documented and justified.

### 2.5 Process execution and external input

Scope:

* Applies to shell commands, subprocess creation, and interactions with untrusted input.

Required rules:

<!-- style: id=2.5.no-shell-interpolation -->
* External command execution MUST avoid shell interpolation when direct argument passing is
  possible.
<!-- style:
id=2.5.validate-untrusted-input
mode=review_only
reason=Input-validation sufficiency depends on cross-boundary data flow.
-->
* Validate and constrain untrusted input before using it in file paths, commands, SQL, or network
  requests.
<!-- style:
id=2.5.secure-temp-files
mode=review_only
reason=Sensitive-temp-file handling depends on whether the data is actually sensitive.
-->
* Temporary files that may hold sensitive data MUST use least-privilege permissions and MUST be
  removed after use.

## 3. Go style

### 3.1 Error handling

Scope:

* Applies to domain, adapter, and application code paths.

Required rules:

<!-- style: id=3.1.domain-errors-location -->
* Domain-specific errors live in `internal/core/domain/errors.go`.
<!-- style: id=3.1.adapters-wrap-with-cause -->
* Adapters MUST wrap low-level errors with context and preserve the cause using `%w`.
<!-- style:
id=3.1.usecases-translate-errors
mode=review_only
reason=Error-translation boundaries still need semantic review.
-->
* Use cases MUST translate infrastructure errors into domain errors before returning them out of
  `core`.
* Error context strings:
  <!-- style: id=3.1.error-context-lowercase -->
  * MUST be lower-case.
  <!-- style: id=3.1.error-context-no-punctuation -->
  * MUST NOT end with punctuation.
  <!-- style: id=3.1.error-context-no-secrets -->
  * MUST NOT include secrets.
<!-- style:
id=3.1.no-string-error-comparison
mode=review_only
reason=String-based error assertions still need repo-wide test cleanup before strict linting.
-->
* Do not compare errors by string content.
<!-- style:
id=3.1.use-errors-is-as
mode=review_only
reason=Equality and assertion checks still need repo-wide cleanup before strict linting.
-->
* Use `errors.Is` and `errors.As` for wrapped errors.
<!-- style:
id=3.1.no-adapter-errors-in-core
mode=review_only
reason=Adapter-error leakage is architectural and often spans type boundaries.
-->
* Do not expose adapter or transport error types from `internal/core`.
<!-- style:
id=3.1.limit-sentinel-errors
mode=review_only
reason=Whether a sentinel is justified remains architectural judgement.
-->
* Use sentinel errors sparingly. Prefer typed or contextual errors when callers need branching.
<!-- style:
id=3.1.no-panic-outside-startup-or-tests
mode=review_only
reason=Allowed impossible-invariant panics still require domain-aware review.
-->
* `panic` is forbidden outside process start-up, impossible invariants, or tests.

Error naming:

<!-- style: id=3.1.sentinel-errors-use-err-prefix -->
* Sentinel errors MUST use an `Err` prefix: `ErrNotFound`, `ErrTimeout`.
<!-- style: id=3.1.error-types-use-error-suffix -->
* Error types MUST use an `Error` suffix: `ValidationError`, `CommandError`.
<!-- style: id=3.1.no-mixed-error-naming -->
* Do not mix the two conventions: `NotFoundError` is not a valid sentinel name, and `ErrValidation`
  is not a valid type name.

Automation guidance:

* Use AST-based checks for lower-case error context, punctuation, secret leakage, sentinel location,
  and adapter `%w` wrapping.

Example:

```go
if err != nil {
    return fmt.Errorf("load profile: %w", err)
}
```

### 3.2 Context, resources, and concurrency

Scope:

* Applies to Go code that uses contexts, I/O resources, goroutines, channels, and shared state.

Required rules:

<!-- style:
id=3.2.pass-context-through
mode=review_only
reason=Correct context propagation across call graphs still needs review.
-->
* Functions that accept `context.Context` MUST pass it through to downstream calls where relevant.
<!-- style: id=3.2.no-context-on-structs -->
* Contexts MUST NOT be stored on structs.
<!-- style: id=3.2.cancel-func-called -->
* Cancel functions returned by `context.WithCancel`, `context.WithTimeout`, or
  `context.WithDeadline` MUST be called on all paths.
<!-- style: id=3.2.close-closers -->
* `io.Closer` values such as files and HTTP bodies MUST be closed.
<!-- style: id=3.2.justify-ignored-close-errors -->
* Ignoring a close error requires an inline comment that explains why it is safe.
<!-- style: id=3.2.explicit-network-timeouts -->
* Network clients MUST use explicit timeouts.
<!-- style:
id=3.2.goroutine-ownership
mode=review_only
reason=Goroutine ownership is behavioural and cross-cutting.
-->
* Do not start goroutines without a defined ownership and shutdown path.
<!-- style:
id=3.2.goroutine-lifecycle
mode=review_only
reason=Goroutine lifecycle reasoning is behavioural and cross-cutting.
-->
* Goroutines MUST NOT outlive the context or lifecycle that created them unless explicitly
  documented.
<!-- style:
id=3.2.synchronise-shared-state
mode=review_only
reason=Shared-state correctness is behavioural rather than purely syntactic.
-->
* Protect shared mutable state with synchronisation.
<!-- style:
id=3.2.single-channel-closer
mode=review_only
reason=Single-closer discipline is behavioural and often cross-function.
-->
* Do not close a channel from multiple senders.
<!-- style: id=3.2.rebind-loop-variables -->
* Loop variables captured by closures MUST be rebound intentionally.

Automation guidance:

* Use `govet`, `staticcheck`, `errcheck`, and targeted custom checks where practical.

### 3.3 Identifier naming

Scope:

* Applies to Go identifiers, types, signatures, and receivers.

Required rules:

Storage naming:

<!-- style: id=3.3.use-repository-name -->
* Use `Repository` for storage interfaces and implementations.
<!-- style: id=3.3.no-store-persistence-types -->
* Do not use `Store` in type names, constructor names, field names, or interface names that
  represent persistence responsibilities.
<!-- style: id=3.3.store-verbs-allowed -->
* The prohibition on `Store` does not apply to verbs such as `StoreMessage`.
<!-- style: id=3.3.prefer-repository-variable-names -->
* Prefer `xRepository` for variables and fields.
<!-- style: id=3.3.prefer-full-identifier-suffixes -->
* Use full forms (`Repository`, `Service`, `Config`, `Manager`) for identifier suffixes, not
  shorthands (`Repo`, `Svc`, `Cfg`, `Mgr`).

Business-logic naming:

<!-- style:
id=3.3.use-service-name
mode=review_only
reason=Workflow-type vocabulary consistency still needs architectural review.
-->
* Use `Service` for domain and application workflow types.
<!-- style:
id=3.3.no-usecase-name
mode=review_only
reason=Legacy workflow vocabulary cleanup still needs repo-wide refactoring judgement.
-->
* Avoid `UseCase` in new code; keep the vocabulary consistent across packages. Variable names:

<!-- style:
id=3.3.boolean-variables-use-state-names
mode=review_only
reason=Whether a bool field needs an is/has/can prefix still needs local readability review.
-->
* Boolean variables and struct fields SHOULD use `is`/`has`/`can`/ `enabled` prefixes or past-
  participle state (`timedOut`, `truncated`, `cancelled`) for clarity. Avoid `flag` and `XFlag`
  suffixes; the bool type already conveys "flag". Short-scope bools (`ok`, `found`, `seen`, `valid`)
  are acceptable when context makes meaning obvious.

<!-- style:
id=3.3.prefer-descriptive-names
mode=review_only
reason=Descriptiveness of names still needs human judgement.
-->
* Prefer descriptive names over abbreviations.
<!-- style: id=3.3.limit-single-letter-names -->
* Single-letter names are acceptable only for loop indices (`i`, `j`, `k`) and method receivers.
* Standard names are allowed:
  <!-- style:
  id=3.3.standard-name-ctx
  mode=review_only
  reason=Signature-name normalisation still needs repo-wide migration planning.
  -->
  * `ctx` for `context.Context`
  <!-- style:
  id=3.3.standard-name-err
  mode=review_only
  reason=Error-name normalisation still needs repo-wide migration planning.
  -->
  * `err` for `error`
<!-- style:
id=3.3.expand-common-abbreviations
mode=review_only
reason=Whether an abbreviation is too terse still needs local naming judgement.
-->
* Avoid abbreviations like `u` for user, `p` for profile, and `f` for file.
<!-- style:
id=3.3.prefer-full-variable-words
mode=review_only
reason=Whether a longer name improves clarity is still contextual.
-->
* Expand them to `user`, `profile`, and `file`.
<!-- style:
id=3.3.name-values-by-role
mode=review_only
reason=The cleanest name still depends on local scope, API shape, and readability tradeoffs.
-->
* Prefer the simplest name that remains immediately clear in local context.
<!-- style:
id=3.3.avoid-overqualified-local-names
mode=review_only
reason=Whether surrounding context supplies enough meaning depends on local scope.
-->
* Avoid over-qualified local variable names when the function, loop, type, or field being accessed
  already supplies the qualifier. Prefer `root`, `path`, or `reference` over names such as
  `repositoryRoot`, `markerPath`, or `configReference` when the shorter name remains immediately
  clear.
<!-- style:
id=3.3.avoid-type-stutter
mode=review_only
reason=Whether a type word is redundant depends on local context.
-->
* Avoid type-stutter and roleless suffixes in local variable names. Do not add words such as
  `Value`, `String`, `Rune`, `Slice`, `Map`, `List`, or `Data` when the type is already obvious from
  local context. Prefer the value's role, such as `character`, `line`, `field`, `target`, or
  `requirement`.
<!-- style:
id=3.3.avoid-acquisition-history-names
mode=review_only
reason=Whether a lighter name stays clear, or a more precise name
genuinely helps, still needs local judgement.
-->
* Prefer lighter, more natural names over heavier or more technical ones when both are equally easy
  to understand.
* Acquisition-history names such as `loadedX`, `parsedX`, or `resolvedX` are acceptable when they
  remain the clearest option.
* Do not rename values solely to make them more systematic or more technically precise if the rename
  does not also make the code read better.
<!-- style:
id=3.3.avoid-opaque-abbreviations
mode=review_only
reason=Opaque abbreviations still need domain-aware review.
-->
* Avoid opaque abbreviations.

Interface naming:

<!-- style: id=3.3.no-i-prefix-on-interfaces -->
* Interface names MUST NOT use an `I` prefix (`IReader`, `IService`). That is a Java convention, not
  Go.
<!-- style:
id=3.3.no-role-suffix-interfaces
mode=review_only
reason=Whether a roleless interface name is justified still needs contextual review.
-->
* Interface names MUST NOT use roleless suffixes such as `Manager`, `Service`, `Helper`, or `Util`.
  Name interfaces after behaviour or a domain concept. Prefer the `-er` suffix for single-method
  behavioural interfaces (`Reader`, `Formatter`).

Public API naming:

<!-- style: id=3.3.no-get-prefix-on-getters -->
* Getters MUST NOT use a `Get` prefix: `obj.Owner()`, not `obj.GetOwner()`.

<!-- style:
id=3.3.keep-exported-api-families-coherent
mode=review_only
reason=Whether an exported variant earns its place depends on caller needs and package scope.
-->
* Keep exported function families minimal and coherent. Do not add an exported variant such as
  `LoadFile` beside `Load` and `Parse` unless it represents a caller-relevant mode, not just a
  convenience wrapper around reading input and calling the source parser.

Helper naming:

<!-- style: id=3.3.boolean-helpers-use-predicate-names -->
* Boolean functions that answer a yes-or-no question MUST be named as predicates:
  `IsX`/`HasX`/`SupportsX`/`MatchesX` when exported, and `isX`/`hasX`/`supportsX`/`matchesX` when
  unexported. Stdlib-style `ValidX`/`CanX` exported forms are not used; the predicate form is
  preferred for uniformity and machine-checkability.
<!-- style:
id=3.3.matches-for-pure-comparisons
mode=review_only
reason=Whether a helper is a pure comparison or an external probe still needs local judgement.
-->
* Reserve `matchesX` names for pure comparisons that do not probe filesystem state, processes, the
  network, or other external state.
<!-- style:
id=3.3.probe-helpers-use-action-verbs
mode=review_only
reason=Probe helpers can still need contextual naming review.
-->
* Use action verbs such as `detect`, `inspect`, `lookup`, or `read` for helpers that probe external
  state, even when they ultimately support a boolean decision.
<!-- style:
id=3.3.keep-helper-verb-families-consistent
mode=review_only
reason=Verb-family consistency still depends on the subsystem's vocabulary and boundaries.
-->
* Use the same verb family for helpers that serve the same role within a package or subsystem.
<!-- style:
id=3.3.use-verbs-for-action-helpers
mode=review_only
reason=Some helpers intentionally name a role or strategy rather than a direct action.
-->
* Prefer verb-led names for helpers that perform work. Reserve noun-led names for types, concepts,
  constructors, and values that represent roles or data.
<!-- style:
id=3.3.reserve-common-helper-verbs
mode=review_only
reason=Verb semantics still need local review when domain language is stronger than generic
  terminology.
-->
* Prefer these verb meanings unless a stronger domain term exists:
  * `parse` for syntactic conversion from text or tokens
  * `read` for filesystem-backed or process-backed input
  * `write` for persisted output
  * `collect` for aggregation across multiple values or files
  * `build` for assembling derived values from already available inputs
  * `check` for rule evaluation or validation decisions
  * `run` for executing a tool, command, or workflow
<!-- style:
id=3.3.avoid-helper-synonyms-without-difference
mode=review_only
reason=Whether two helper verbs are genuinely distinct still needs API-level judgement.
-->
* Do not alternate between helper verbs such as `build`, `create`, `assemble`, `check`, `validate`,
  or `verify` unless the names express a real semantic difference.

Receiver names:

<!-- style:
id=3.3.short-receiver-names
mode=review_only
reason=Receiver brevity and local readability still need contextual review.
-->
* Method receiver names MUST be short.
<!-- style:
id=3.3.receiver-abbreviates-type
mode=review_only
reason=Acceptable receiver abbreviations still need contextual review.
-->
* Receiver names MUST abbreviate the receiver type.
<!-- style:
id=3.3.receiver-consistency
mode=review_only
reason=Receiver consistency across large types still needs repo-wide cleanup.
-->
* Receiver names MUST be consistent across all methods on that type.

Go identifier forms:

<!-- style: id=3.3.acronyms-upper-in-mixedcaps -->
* Acronyms and initialisms MUST be fully upper-cased within MixedCaps identifiers: `UserID`,
  `ServeHTTP`, `ParseURL`. Do not write `UserId`, `ServeHttp`, or `ParseUrl`.

<!-- style:
id=3.3.no-underscores-in-identifiers
mode=review_only
reason=Interop and test-name exceptions still need contextual review.
-->
* Go identifier names MUST NOT contain underscores except where the Go toolchain or interop
  conventions require them.
<!-- style:
id=3.3.underscore-exceptions
mode=review_only
reason=Underscore exceptions remain contextual and tied to framework conventions.
-->
* Allowed exceptions include test names, generated-code imports, and rare low-level interop cases.

Automation guidance:

* Use AST checks for single-letter variable names.
* Use text checks for identifier suffixes (`Repo`, `Svc`, `Cfg`, `Mgr`) and type suffixes (`Store`),
  mapped to their preferred full forms (`Repository`, `Service`, `Config`, `Manager`).

Correct:

```go
func LoadIdentity(
    ctx context.Context,
    id domain.IdentityID,
    passphrase string,
) (identity domain.Identity, err error) {
    identity, err = ...
    if err != nil {
        return domain.Identity{}, err
    }
    return identity, nil
}

func DoSomething(ctx context.Context, data string, passphrase string) (err error) {
    err = ...
    if err != nil {
        return err
    }
    return nil
}
```

```go
func Send(
    ctx context.Context,
    from domain.Username,
    to domain.Username,
    message string,
) (err error)
```

```go
identityID, err := domain.ParseIdentityID(rawIdentityID)
if err != nil {
    return err
}

peer, err := domain.ParseUsername(rawPeer)
if err != nil {
    return err
}
```

```go
func (s *Service) LoadProfile(
    username domain.Username,
) (profile domain.Profile, err error) {
    profile, found, err := s.profileRepository.Load(username)
    if err != nil {
        return domain.Profile{}, err
    }
    _ = found
    return profile, nil
}

file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()

for i := range items {
    // ok: loop index
}
```

Incorrect:

```go
func LoadIdentity(
    ctx context.Context,
    id domain.IdentityID,
    passphrase string,
) (domain.Identity, error)

func LoadIdentity(...) (identity domain.Identity, err error) {
    ...
    return
}
```

```go
func Send(ctx context.Context, from, to domain.Username, message string) error
```

```go
identityID := domain.IdentityID(rawIdentityID)
peer := domain.Username(rawPeer)
conversationID := domain.ConversationID(rawConversationID)
```

```go
func (s *Service) LoadProfile(u domain.Username) (profile domain.Profile, err error) {
    p, f, e := s.profileRepository.Load(u)
    if e != nil {
        return domain.Profile{}, e
    }
    _ = p
    _ = f
    return profile, nil
}
```

### 3.4 Function signatures

Scope:

* Applies to function and method return values and parameters.

Required rules:

Return values:

<!-- style: id=3.4.named-return-values -->
* All Go functions and methods MUST use named return values.
<!-- style: id=3.4.named-returns-in-all-signatures -->
* This rule applies to declared functions, methods, interface methods, and function types.
<!-- style:
id=3.4.named-anon-returns-when-practical
mode=review_only
reason=Whether named returns improve local anonymous-function clarity is contextual.
-->
* Anonymous functions SHOULD use named return values when practical.
<!-- style:
id=3.4.meaningful-return-names
mode=review_only
reason=Meaningfulness of names still needs human judgement.
-->
* Return names MUST be meaningful and descriptive.
<!-- style: id=3.4.no-placeholder-return-names -->
* Placeholder return names such as `result0`, `result1`, and similar are not allowed.
<!-- style: id=3.4.no-naked-returns -->
* Do not use naked returns in functions that return values.
<!-- style: id=3.4.explicit-return-values -->
* Always return explicit values.
<!-- style: id=3.4.naked-return-void-only -->
* Naked `return` is allowed only in functions with no return values.

Parameter types:

<!-- style: id=3.4.no-type-elision -->
* Avoid type elision in function signatures.
<!-- style: id=3.4.explicit-parameter-types -->
* Each parameter SHOULD have its own type declaration for clarity.
<!-- style: id=3.4.no-direct-domain-casts -->
* Do not directly cast raw strings to key domain identifier aliases: `domain.Username`,
  `domain.ConversationID`, `domain.IdentityID`.
<!-- style: id=3.4.use-domain-parsers -->
* Outside `internal/core/domain`, use parser or constructor helpers instead: `ParseUsername`,
  `ParseConversationID`, `ConversationIDFromUsername`, `ParseIdentityID`.
<!-- style: id=3.4.domain-cast-rule-in-tests -->
* This rule applies to both production and test Go code.

Automation guidance:

* Use AST checks for named returns, naked returns, type elision, and domain identifier constructor
  usage.

### 3.5 Ordering within types

Method ordering:

Scope:

* Applies to ports interfaces, mocks, and implementations.

Required rules:

<!-- style: id=3.5.crudl-order -->
* Use CRUD-L order: Create -> Read -> Update -> Delete -> List.
<!-- style:
id=3.5.crudl-section-headers
mode=review_only
reason=Whether a type has enough method groups to justify CRUD-L headers is contextual.
-->
* Use section headers (`Create`, `Read`, `Update`, `Delete`, `List`) in types with multiple groups.
<!-- style:
id=3.5.omit-empty-crudl-groups
mode=review_only
reason=Whether a CRUD-L group is meaningfully empty still depends on interface design
  intent.
-->
* If a group has no methods, omit the group, but keep remaining groups in CRUD-L order.
<!-- style: id=3.5.mocks-match-interface-order -->
* Mocks MUST match interface order exactly.
<!-- style: id=3.5.implementations-match-interface-order -->
* Implementations MUST match interface order.

Automation guidance:

* Use AST checks to enforce interface CRUD-L ordering and mock ordering parity.

Example:

```go
type MessageRepository interface {
    /* ------------------------------------------- Create --------------------------------------- */

    SaveMessage(ctx context.Context, ...) (err error)

    /* -------------------------------------------- Read ---------------------------------------- */

    ListMessages(ctx context.Context, ...) (messages []Message, err error)
    ListConversations(ctx context.Context) (ids []ConversationID, err error)
}
```

Constructor ordering:

Scope:

* Applies to constructors, for example `New` functions.

Required rules:

<!-- style: id=3.5.constructor-category-order -->
* Constructors MUST order dependencies by category:

1. Repositories (persistence) Canonical order:
   1. `identityRepository` (`IdentityRepository`)
   2. `preKeyRepository` (`PreKeyRepository`)
   3. `preKeyBundleRepository` (`PreKeyBundleRepository`)
   4. `sessionRepository` (`SessionRepository`)
   5. `ratchetRepository` (`RatchetRepository`)
   6. `accountRepository` (`AccountRepository`)

2. Services (domain dependencies) Canonical order:
   1. `identityService`
   2. `preKeyService`
   3. `sessionService`
   4. `messageService`

3. Adapters (clients or factories) Canonical order:
   1. `relayClient`
   2. `relayClientFactory`

4. Configuration or identity Canonical order:
   1. `serverURL` or `relayURL`
   2. `identityID`
   3. `timeout`

5. Secrets, always last

Rules for dependencies not in the canonical list:

<!-- style:
id=3.5.group-unlisted-dependencies
mode=review_only
reason=How unlisted constructor dependencies group is still contextual.
-->
* Group them by the categories above.
<!-- style:
id=3.5.place-unlisted-after-canonical
mode=review_only
reason=Placement of unlisted constructor dependencies remains contextual.
-->
* Place them after listed dependencies within the same category.
<!-- style:
id=3.5.sort-unlisted-by-name
mode=review_only
reason=Alphabetising unlisted constructor dependencies still depends on the chosen grouping
  strategy.
-->
* Sort them alphabetically by parameter name within that category.

Additional API design rules:

<!-- style:
id=3.5.prefer-options-pattern
mode=review_only
reason=Constructor design still needs API-level review.
-->
* Constructors with multiple optional arguments SHOULD prefer an options pattern.
<!-- style:
id=3.5.consider-functional-options
mode=review_only
reason=Constructor design still needs API-level review.
-->
* Consider functional options for exported APIs that are expected to grow.
<!-- style:
id=3.5.no-positional-optional-args
mode=review_only
reason=Whether an argument is truly optional still depends on API design intent.
-->
* Do not add positional optional arguments to exported constructors.

Automation guidance:

* Use AST checks to enforce constructor dependency category order.

Example:

```go
func New(
    identityRepository ports.IdentityRepository,
    sessionService ports.SessionService,
    relayClientFactory ports.RelayClientFactory,
    serverURL string,
) (service *Service)
```

### 3.6 Imports

Scope:

* Applies to all Go files.

Required rules:

<!-- style: id=3.6.use-gofmt -->
* Use `gofmt -s` for formatting.
<!-- style: id=3.6.use-goimports -->
* Use `goimports` to manage imports and grouping.
<!-- style: id=3.6.import-group-order -->
* Keep import groups in this order:
  1. Standard library
  2. Third party
  3. Local project packages
<!-- style: id=3.6.no-blank-imports -->
* Blank imports (`import _`) are forbidden outside `main` packages and tests, except for documented
  toolchain exceptions such as `embed`.
<!-- style:
id=3.6.alias-only-when-needed
mode=review_only
reason=Whether an alias is genuinely needed still depends on the local naming context.
-->
* Import aliases are forbidden unless:
  * the alias resolves a real naming conflict in the importing file, or
  * a third-party package's natural import name is unclear or awkward enough that an explicit alias
    improves readability.
<!-- style:
id=3.6.prefer-natural-package-names
mode=review_only
reason=Whether repeated aliases indicate a package-naming problem
still needs subsystem-level judgement.
-->
* Package names SHOULD read naturally at import sites.
<!-- style:
id=3.6.reconsider-repeated-aliases
mode=review_only
reason=Deciding whether repeated aliases justify an API or package
rename still needs design judgement.
-->
* If many files need the same non-conflicting alias for one package, reconsider the package name or
  API surface instead of normalising the alias.

Automation guidance:

* Use formatting and linting tools to enforce import order.

Example:

```go
import (
    "context"
    "fmt"

    "github.com/google/uuid"

    "your/module/internal/core/domain"
)
```

<!-- style: id=3.6.no-package-type-stutter -->
* Exported type names MUST NOT repeat the package name: `style.StyleRule` stutters; use
  `style.Rule`. The mild form `pkg.Pkg` (for example `time.Time`) is tolerated only for foundational
  types where the repetition reads naturally; otherwise rename the type or the package. The same
  rule applies to exported functions: `style.NewStyle` stutters; `style.New` does not.

### 3.7 Parameter ordering

Scope:

* Applies to function and method signatures with multiple parameters.

Required rules:

<!-- style: id=3.7.ctx-first -->
1. Context (`ctx context.Context`), always first if present.
<!-- style:
id=3.7.subject-before-options
mode=review_only
reason=Subject-versus-options ordering depends on API shape and call-site clarity.
-->
1. Subject identifier, the noun being acted upon
<!-- style:
id=3.7.options-before-payload
mode=review_only
reason=Options-versus-payload ordering depends on API shape and call-site clarity.
-->
1. Non-secret auth, config, or options
<!-- style:
id=3.7.payload-before-secrets
mode=review_only
reason=Payload-versus-secrets ordering depends on API shape and call-site clarity.
-->
1. Data or payload
<!-- style: id=3.7.secrets-last -->
1. Secrets, always last

Examples:

* Subject identifier: `identityID`, `username`, `peer`, `conversationID`
* Non-secret config or options: `relayURL`, `serverURL`, `timeout`
* Payload: `message`, `count`, `metadata`
* Secrets: `passphrase`, `privateKey`, `token`, `seed`

Automation guidance:

* Use AST checks for `ctx` first and secrets last ordering.

Example:

```go
func (u *SendMessageUseCase) Execute(
    ctx context.Context,
    conversationID string,
    relayURL string,
    message domain.Message,
    passphrase string,
) (err error)
```

### 3.8 Comments

Scope:

* Applies to Go doc comments, block comments, and inline comments.

Required rules:

Doc comments (exported identifiers):

<!-- style: id=3.8.doc-comments-start-name -->
* MUST start with the identifier name.
<!-- style: id=3.8.doc-comments-full-sentences -->
* Use full sentences and end with a period.
<!-- style:
id=3.8.doc-comments-describe-what
mode=review_only
reason=Distinguishing behaviour from existence-justification is a readability judgement.
-->
* Describe what the identifier does, with its mechanism. A second sentence is
  allowed only to explain a non-obvious constraint or parameter - never to
  justify why the identifier exists. Existence rationale belongs in commit
  messages, not doc comments.

Block comments (multi-line explanations and single-line comments that annotate a nearby block):

<!-- style:
id=3.8.block-comments-full-sentences
mode=review_only
reason=Banner comments and block-local annotations still need contextual review.
-->
* Use full sentences and end with a period.

Inline comments (trailing comments):

<!-- style:
id=3.8.inline-comments-short
mode=review_only
reason=Comment brevity is a readability judgement rather than a syntax check.
-->
* Keep short.
<!-- style: id=3.8.inline-comments-lowercase -->
* Start with a lower-case letter.
<!-- style: id=3.8.inline-comments-no-period -->
* Do not end with a period.
<!-- style:
id=3.8.inline-comments-explain-why
mode=review_only
reason=Comment quality is better enforced in review than lint.
-->
* Prefer explaining why over restating what.

Automation guidance:

* Use AST-based checks for inline trailing comment case and punctuation.

Examples:

```go
// downloadFile fetches a URL and writes it to destination via an atomic
// temp-file rename. The download is capped at limit bytes to prevent
// unbounded memory or disk usage.
//
// Preferred: both sentences describe behaviour and mechanism. The second
// explains a constraint, not why the function exists.

// hasPinnedLocalTool reports whether a tool matching the pinned version is
// already installed at the given path.
// This makes re-running install idempotent.
//
// Avoid: the second sentence justifies existence, not behaviour. It belongs
// in a commit message, not a doc comment.

if len(key) == 0 {
    return domain.ErrMissingKey
}

nonce := make([]byte, 24) // generated per message
```

### 3.9 File structure and package lifecycle

Scope:

* Applies to Go files.

Required rules:

<!-- style:
id=3.9.order-elements-consistently
mode=review_only
reason=Consistency across file elements is broader than the objective ordering checks.
-->
* Order elements consistently.
<!-- style: id=3.9.file-order -->
* Use this order:
  1. Package declaration and imports
  2. Constants Header: `/* ... Constants ... */`
  3. Errors (sentinel errors, if any)
  4. Types Header: `/* ... Types ... */` Exported types first, then internal types.
  5. Constructor Header: `/* ... Constructor ... */`
  6. Methods in CRUD-L order with section headers
  7. Helpers Header: `/* ... Helpers ... */`
  8. Compile-time assertions Header: `/* ... Assert ... */`
<!-- style:
id=3.9.reader-journey-order
mode=review_only
reason=Reader journey depends on the file's role and cannot be reduced to declaration syntax.
-->
* Order each file by reader journey: public entrypoint, main flow, specialised branches, narrow
  helpers. Use declaration category as the primary order only in mixed structural model files.
<!-- style:
id=3.9.role-specific-file-grammar
mode=review_only
reason=Role-specific ordering grammars require package-aware judgement.
-->
* Files with the same platform role SHOULD follow the same internal ordering grammar: CLI commands
  use `runX`, option parsing, flag construction, selection/execution, then rendering/status helpers.
  Profile schema, conversion, and validation files follow `quill.toml` policy order. Rule scanners
  use exported check, file collection, scan/traversal, diagnostic construction, parsing/extraction
  helpers, then predicates. Drivers use driver entrypoint, dispatch, concrete Driver, then command
  helpers. Runtime installers use orchestration, inspection, preparation and validation, install
  execution, then low-level filesystem helpers. Report files use public writer, text rendering, JSON
  conversion, then formatting helpers.
<!-- style:
id=3.9.split-by-responsibility
mode=review_only
reason=Ciphera application code does not yet satisfy automated file-shape thresholds.
-->
* Split files by domain responsibility, not by habit. Prefer a calm tree over file confetti.
<!-- style:
id=3.9.merge-tiny-glue-files
mode=review_only
reason=Ciphera application code does not yet satisfy automated file-shape thresholds.
-->
* Merge tiny glue files when they only hold one narrow helper, alias, or constant and make readers
  jump between files without reducing complexity.
<!-- style:
id=3.9.avoid-vague-file-names
mode=review_only
reason=Ciphera application code does not yet satisfy automated file-shape thresholds.
-->
* Avoid vague filenames such as `types.go`, `helpers.go`, `model.go`, and `checks.go` unless the
  file is genuinely package-wide.
<!-- style:
id=3.9.long-files-need-justification
mode=review_only
reason=Ciphera application code does not yet satisfy automated file-shape thresholds.
-->
* Handwritten Go files over 300 lines and non-test functions over 80 lines need a clear structural
  reason or should be split.
<!-- style:
id=3.9.names-add-information
mode=review_only
reason=Naming quality depends on local context and reader cost.
-->
* Names should add information not already obvious from the type or immediate context. Avoid type-
  restating names such as `runeValue` when `rune` or a domain noun is clearer.
<!-- style: id=3.9.scanner-entrypoint-first -->
* First-party scanner files MUST place exported `Check...` entrypoints before unexported scan,
  diagnostic, parsing, and predicate helpers.

Additional required rules:

<!-- style:
id=3.9.no-package-mutable-vars
mode=review_only
reason=Compatibility and test-seam exceptions need package-level review.
-->
* Package-level mutable variables are forbidden unless required for compatibility or test seams.
<!-- style:
id=3.9.init-requires-justification
mode=review_only
reason=init() justifications need package-level review.
-->
* `init()` functions are discouraged and require justification.
<!-- style:
id=3.9.init-must-be-deterministic
mode=review_only
reason=init() determinism requires behavioural review beyond syntax.
-->
* If `init()` is required, it MUST be deterministic and MUST NOT depend on environment, filesystem
  state, process arguments, network state, or ordering side effects from other files.
<!-- style:
id=3.9.group-methods-by-receiver
mode=review_only
reason=Receiver-grouping tradeoffs still depend on file readability.
-->
* Group methods by receiver where practical.
<!-- style:
id=3.9.keep-constructor-adjacent
mode=review_only
reason=Constructor adjacency still depends on the file's broader layout.
-->
* Keep constructors adjacent to the primary exported type they construct.

Automation guidance:

* Enforce broad reader-journey order through code review.
* Use AST checks for objective file-ordering cases: declaration group order, scanner entrypoint
  order, and test-helper placement.
* Use recommendation checks for file-shape signals after the codebase is clean enough that strict
  recommendations do not become noisy.

Example:

```go
package storage

/* ------------------------------------------ Constants ---------------------------------------- */

const maxRetries = 3

/* -------------------------------------------- Types ------------------------------------------ */

type Repository struct {
    ...
}

/* ----------------------------------------- Constructor --------------------------------------- */

func NewRepository(...) (repository *Repository) { ... }

/* ------------------------------------------- Create ------------------------------------------ */

func (r *Repository) Save(...) (err error) { ... }

/* -------------------------------------------- Read ------------------------------------------- */

func (r *Repository) Load(...) (value ..., err error) { ... }

/* ------------------------------------------- Helpers ----------------------------------------- */

func (r *Repository) ensureDir(...) (err error) { ... }

/* -------------------------------------------- Assert ----------------------------------------- */

var _ ports.Repository = (*Repository)(nil)
```

### 3.10 Data and interface usage

Scope:

* Applies to Go slices, maps, struct literals, interfaces, and zero-value patterns.

Required rules:

<!-- style: id=3.10.named-struct-literals -->
* Use field names in struct literals by default.
<!-- style:
id=3.10.limit-positional-struct-literals
mode=review_only
reason=Whether a literal is tuple-like or trivial still needs local judgement.
-->
* Positional struct literals are allowed only for very small local test tables or where the type is
  intentionally tuple-like.
<!-- style:
id=3.10.nil-default-empty-slices
mode=review_only
reason=Nil-versus-empty decisions still depend on external contract semantics.
-->
* Treat `nil` as the default zero-length slice unless a non-`nil` empty slice is required by an
  external contract.
<!-- style: id=3.10.use-len-for-slice-emptiness -->
* Test slice emptiness with `len(slice) == 0`.
<!-- style:
id=3.10.prefer-usable-zero-values
mode=review_only
reason=Usable-zero-value tradeoffs still depend on API intent and mutation semantics.
-->
* Prefer usable zero values over unnecessary allocation.
<!-- style: id=3.10.no-pointers-to-interfaces -->
* Do not use pointers to interfaces.
<!-- style: id=3.10.pass-interface-values-directly -->
* Pass interface values directly.

Automation guidance:

* Enforce through code review and targeted linting where practical.

## 4. Bash style

### 4.1 Naming and structure

Scope:

* Applies to all `.sh` files in the repository.

Required rules:

<!-- style: id=4.1.shebang -->
* Scripts MUST start with `#!/bin/bash`.
<!-- style: id=4.1.strict-mode -->
* Scripts MUST include `set -euo pipefail`.
<!-- style: id=4.1.tabs-not-spaces -->
* Leading indentation MUST use tabs, not spaces.
<!-- style: id=4.1.no-trailing-whitespace -->
* Avoid trailing whitespace.
<!-- style: id=4.1.unix-line-endings -->
* Use Unix line endings.
<!-- style: id=4.1.follow-language-guidance -->
* Comments and user-facing output MUST follow section 1.3 language guidance.
<!-- style: id=4.1.prefer-descriptive-shell-names -->
* Prefer descriptive names over abbreviations.
<!-- style: id=4.1.lowercase-shell-functions -->
* Function names and non-exported variable names SHOULD use lower-case with underscores.
<!-- style: id=4.1.uppercase-shell-constants -->
* Constants and exported environment variables SHOULD use upper-case with underscores.
<!-- style: id=4.1.main-at-bottom -->
* Non-trivial scripts MUST define `main()` as the bottom-most function and end with `main "$@"`.
<!-- style:
id=4.1.shell-crudl-order
mode=review_only
reason=Not every shell script models a resource, but CRUD-L ordering improves scanability when
  it does.
-->
* When a shell script exposes multiple handlers for one resource, order the handlers in CRUD-L order
  where the semantics fit.
<!-- style:
id=4.1.shell-handler-family-names
mode=review_only
reason=Handler-family naming still depends on whether the script models one resource clearly.
-->
* When a shell script manages one resource, use a consistent action-plus-resource naming family for
  handlers and probes.
<!-- style:
id=4.1.shell-predicates-vs-actions
mode=review_only
reason=Probe-versus-action naming still needs local judgement in shell code.
-->
* Name state probes as predicates such as `is_x`, `has_x`, or `matches_x`, and name mutating
  handlers with action verbs such as `create_x`, `read_x`, `update_x`, `delete_x`, `list_x`,
  `start_x`, or `stop_x`.

Automation guidance:

* Use Bash-specific style checks for shebang, strict mode, tabs, whitespace, and line endings.
* Use `shfmt` to enforce Bash formatting.

### 4.2 Shell safety and correctness

Scope:

* Applies to shell commands, substitutions, argument handling, temporary files, and suppressions.

Required rules:

<!-- style: id=4.2.command-substitution -->
* Use `$(...)` for command substitution. Do not use backticks.
<!-- style: id=4.2.prefer-double-brackets -->
* Use `[[ ... ]]` instead of `[ ... ]` where Bash semantics are intended.
<!-- style: id=4.2.quote-expansions -->
* Quote all parameter expansions unless unquoted splitting is explicitly required and documented.
<!-- style:
id=4.2.use-arrays-for-args
mode=review_only
reason=Whether an argument list warrants arrays still depends on shell-script context.
-->
* Use arrays for command argument lists and expand them as `"${array[@]}"`.
<!-- style:
id=4.2.prefer-printf
mode=review_only
reason=printf versus echo remains a local shell-readability judgement when output is
  simple.
-->
* Prefer `printf` over `echo` for user-facing output.
<!-- style: id=4.2.use-read-r -->
* Use `read -r`.
<!-- style:
id=4.2.local-shell-variables
mode=review_only
reason=Function-local scoping in Bash remains contextual without a real shell AST.
-->
* Declare function-local variables with `local`.
<!-- style: id=4.2.avoid-eval -->
* Avoid `eval`.
<!-- style: id=4.2.detect-tools-with-command-v -->
* Use `command -v` to detect dependencies.
<!-- style: id=4.2.use-mktemp -->
* Temporary files and directories MUST be created with `mktemp`.
<!-- style: id=4.2.trap-temp-cleanup -->
* Cleanup for temporary resources MUST use `trap`.
<!-- style: id=4.2.avoid-state-losing-read-loops -->
* Prefer process substitution or `readarray` over `cmd | while read ...` when state must survive
  after the loop.
<!-- style: id=4.2.local-shellcheck-suppressions -->
* ShellCheck suppressions MUST be local, MUST identify the rule when possible, and MUST include a
  short reason.

Automation guidance:

* Use `shellcheck -x` for Bash static analysis.

Example:

```bash
#!/bin/bash
set -euo pipefail

# ------------------------------------------- Helpers -------------------------------------------

render_line() {
	local text="$1"
	printf '%s\n' "$text"
}

main() {
	render_line "ready"
}

main "$@"
```

## 5. Documentation

### 5.1 Audience

Scope:

* Applies to project documentation and developer-facing explanatory content.

Required rules:

<!-- style:
id=5.1.write-for-amateur-security-readers
mode=review_only
reason=Audience fit is a writing judgement rather than a lint rule.
-->
* Write for developers with amateur cybersecurity knowledge.
<!-- style:
id=5.1.explain-security-plainly
mode=review_only
reason=Plain-language quality is a writing judgement rather than a lint rule.
-->
* Explain security concepts plainly and avoid jargon where possible.
<!-- style:
id=5.1.define-jargon-once
mode=review_only
reason=Jargon handling is a documentation judgement rather than a lint rule.
-->
* If jargon is necessary, define it once in simple terms.

### 5.2 Format

Scope:

* Applies to `.md` files, Go comments, and `doc.go` documentation.

Required rules:

<!-- style: id=5.2.use-markdown -->
* Use Markdown for `.md` files.
<!-- style: id=5.2.keep-go-docs-readable -->
* In Go comments and `doc.go`, avoid Markdown syntax that reduces readability.
<!-- style:
id=5.2.full-sentences-and-clear-intent
mode=review_only
reason=Documentation quality is better enforced in review than lint.
-->
* Prefer full sentences and clear intent.
<!-- style:
id=5.2.comments-follow-go-comment-rules
mode=review_only
reason=This is an umbrella documentation rule over both automated and review-driven comment
  guidance.
-->
* Comments MUST follow the rules in section 3.8.
<!-- style:
id=5.2.concise-and-clear
mode=review_only
reason=Documentation quality is better enforced in review than lint.
-->
* Be concise while staying clear for the target audience.

Automation guidance:

* Use `markdownlint` to enforce Markdown structure and formatting.

### 5.3 Markers

Scope:

* Applies to inline maintenance notes in code and documentation.

Required rules:

<!-- style:
id=5.3.todo-fixme-allowed
mode=review_only
reason=This is an allowance policy rather than a standalone lint obligation.
-->
* `TODO:` and `FIXME:` are allowed.
<!-- style: id=5.3.markers-short-and-actionable -->
* Keep them short and actionable.
<!-- style: id=5.3.markers-identify-next-step -->
* Include enough context to identify the next concrete step.

## 6. Tests

### 6.1 Go tests

Scope:

* Applies to `_test.go` files.

Required rules:

<!-- style: id=6.1.helpers-call-helper -->
* Test helpers MUST call `t.Helper()`.
<!-- style: id=6.1.helpers-after-tests -->
* Test helpers SHOULD appear below the test cases in their file.
<!-- style: id=6.1.use-tempdir -->
* Tests MUST use `t.TempDir()` for temporary filesystem state.
<!-- style: id=6.1.use-setenv -->
* Tests MUST use `t.Setenv()` for environment changes.
<!-- style: id=6.1.avoid-arbitrary-sleeps -->
* Time-based tests MUST avoid arbitrary sleeps where a deterministic signal is possible.
<!-- style:
id=6.1.prefer-table-driven-tests
mode=review_only
reason=Test structure preference is better enforced in review than lint.
-->
* Prefer table-driven tests when the test logic is repetitive across cases.
<!-- style:
id=6.1.use-subtests
mode=review_only
reason=Subtest structure preference is better enforced in review than lint.
-->
* Use subtests for named cases.
<!-- style:
id=6.1.one-behaviour-per-table
mode=review_only
reason=Test case granularity is better enforced in review than lint.
-->
* Keep each test table focused on one behaviour.
<!-- style:
id=6.1.relevant-table-fields
mode=review_only
reason=Table field relevance is better enforced in review than lint.
-->
* All fields in a test table SHOULD be relevant to all rows.
<!-- style:
id=6.1.named-fields-in-non-trivial-tables
mode=review_only
reason=Table readability is better enforced in review than lint.
-->
* Prefer named fields in test-case struct literals when the table is non-trivial.
<!-- style:
id=6.1.include-failure-context
mode=review_only
reason=Failure-message quality is better enforced in review than lint.
-->
* Test failures MUST include enough context to diagnose the failing input.

### 6.2 Bash tests

Scope:

* Applies to Bash tests and shell-based test helpers.

Required rules:

<!-- style: id=6.2.cleanup-temp-resources -->
* Tests MUST clean up temporary files and directories.
<!-- style:
id=6.2.avoid-ambient-environment
mode=review_only
reason=Ambient-environment reliance often requires behavioural review across helper layers.
-->
* Tests MUST avoid depending on ambient environment state when practical.
<!-- style:
id=6.2.make-failures-obvious
mode=review_only
reason=Failure-output clarity is better enforced in review than lint.
-->
* Test output SHOULD make the failing command or input obvious.
